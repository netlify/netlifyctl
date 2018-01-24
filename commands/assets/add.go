package assets

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/spf13/cobra"
)

const (
	uploadedAssetState     = "uploaded"
	privateAssetVisibility = "private"
)

type assetsAddCmd struct {
	private bool
}

func setupAddCommand(middlewares []middleware.Middleware) *cobra.Command {
	cmd := &assetsAddCmd{}
	ccmd := &cobra.Command{
		Use:   "add [ASSET PATH 1] [ASSET PATH 2] ...",
		Short: "Add an asset to a site",
		Long:  "Add an asset to a site",
	}
	ccmd.Flags().BoolVarP(&cmd.private, "private", "p", false, "make the asset private")

	return middleware.SetupCommand(ccmd, cmd.AddAsset, middlewares)
}

func (c *assetsAddCmd) AddAsset(ctx context.Context, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing asset paths to upload")
	}

	conf := context.GetSiteConfig(ctx)
	if conf.Settings.ID == "" {
		return errors.New("Failed to load site configuration")
	}

	dups := make(map[string]os.FileInfo)
	for _, arg := range args {
		fi, err := os.Stat(arg)
		if err != nil {
			return fmt.Errorf("unable to access asset with path %s, all assets must exist: %v", arg, err)
		}

		if _, dup := dups[arg]; dup {
			fmt.Printf("[WARNING] duplicated asset %s, will be uploaded only once\n", arg)
		}

		dups[arg] = fi
	}

	params := operations.NewCreateSiteAssetParams().WithSiteID(conf.Settings.ID)
	if c.private {
		visibility := privateAssetVisibility
		params = params.WithVisibility(&visibility)
	}

	for fp, fi := range dups {
		asset, err := c.uploadAsset(ctx, conf.Settings.ID, fp, fi, params)
		if err != nil {
			fmt.Printf("%s upload failed: %v", fp, err)
			break
		}

		fmt.Printf("%s available at %s\n", fp, asset.URL)
	}

	return nil
}

func (c *assetsAddCmd) uploadAsset(ctx context.Context, siteId string, fp string, fi os.FileInfo, baseParams *operations.CreateSiteAssetParams) (*models.Asset, error) {
	client := context.GetClient(ctx)

	body, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	buffer := make([]byte, 512)
	n, err := body.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}
	contentType := http.DetectContentType(buffer[:n])
	body.Seek(0, 0)

	params := baseParams.WithName(fi.Name()).WithSize(fi.Size()).WithContentType(contentType)

	signature, err := client.AddSiteAsset(ctx, params)
	if err != nil {
		return nil, err
	}

	bufferWriter := &bytes.Buffer{}
	writer := multipart.NewWriter(bufferWriter)
	defer writer.Close()

	for key, value := range signature.Form.Fields {
		writer.WriteField(key, value)
	}

	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, body); err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", signature.Form.URL, bufferWriter)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("unexpected error uploading assets to the :cloud: %d - %s", resp.StatusCode, respBody)
	}

	updateParams := operations.NewUpdateSiteAssetParams().WithSiteID(siteId).WithAssetID(signature.Asset.ID).WithState(uploadedAssetState)

	return client.UpdateSiteAsset(ctx, updateParams)
}
