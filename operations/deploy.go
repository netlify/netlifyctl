package operations

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/netlify/netlifyctl/ui"
	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/porcelain"
)

const progressColor = "white"

type DeployObserver struct {
	walkerSpinner *spinner.Spinner
	deltaSpinner  *spinner.Spinner
	uploadSpinner *spinner.Spinner

	total    int
	delta    int
	uploaded int

	uploadedC chan *porcelain.FileBundle
	uploadedD chan bool
}

func NewDeployObserver() *DeployObserver {
	return &DeployObserver{
		walkerSpinner: spinner.New(spinner.CharSets[39], 300*time.Millisecond),
		deltaSpinner:  spinner.New(spinner.CharSets[39], 300*time.Millisecond),
		uploadSpinner: spinner.New(spinner.CharSets[39], 300*time.Millisecond),
		uploadedC:     make(chan *porcelain.FileBundle),
		uploadedD:     make(chan bool),
	}
}

func (o *DeployObserver) OnSetupWalk() error {
	o.walkerSpinner.Prefix = "Counting objects .... "
	o.walkerSpinner.Color(progressColor)
	o.walkerSpinner.Start()
	return nil
}

func (o *DeployObserver) OnSuccessfulStep(*porcelain.FileBundle) error {
	o.total += 1
	o.walkerSpinner.Prefix = fmt.Sprintf("Counting objects: %d ", o.total)
	return nil
}

func (o *DeployObserver) OnSuccessfulWalk(df *models.DeployFiles) error {
	o.walkerSpinner.FinalMSG = fmt.Sprintf("Counting objects: %d total objects  %s\n", o.total, ui.DoneCheck())
	o.walkerSpinner.Stop()
	return nil
}

func (o *DeployObserver) OnFailedWalk() {
	o.walkerSpinner.FinalMSG = fmt.Sprintf("Counting objects  %s\n", ui.ErrorCheck())
	o.walkerSpinner.Stop()
}

func (o *DeployObserver) OnSetupDelta(*models.DeployFiles) error {
	o.deltaSpinner.Prefix = "Resolving deltas .... "
	o.deltaSpinner.Color(progressColor)
	o.deltaSpinner.Start()
	return nil
}

func (o *DeployObserver) OnSuccessfulDelta(df *models.DeployFiles, d *models.Deploy) error {
	o.delta = len(d.Required) + len(d.RequiredFunctions)

	o.deltaSpinner.FinalMSG = fmt.Sprintf("Resolving deltas: %d objects to upload  %s\n", o.delta, ui.DoneCheck())
	o.deltaSpinner.Stop()

	go o.listenUploads()
	o.uploadSpinner.Prefix = fmt.Sprintf("Uploading objects: %d/%d ", o.uploaded, o.delta)
	o.uploadSpinner.Color(progressColor)
	o.uploadSpinner.Start()
	return nil
}

func (o *DeployObserver) OnFailedDelta(*models.DeployFiles) {
	o.deltaSpinner.FinalMSG = fmt.Sprintf("Resolving deltas  %s\n", ui.ErrorCheck())
	o.deltaSpinner.Stop()
}

func (o *DeployObserver) OnSetupUpload(f *porcelain.FileBundle) error {
	return nil
}

func (o *DeployObserver) OnSuccessfulUpload(f *porcelain.FileBundle) error {
	o.uploadedC <- f
	return nil
}

func (o *DeployObserver) OnFailedUpload(*porcelain.FileBundle) {
	o.uploadedD <- true
	o.uploadSpinner.FinalMSG = fmt.Sprintf("Uploading objects  %s\n", ui.ErrorCheck())
	o.uploadSpinner.Stop()
}

func (o *DeployObserver) Finish() {
	o.uploadedD <- true
	o.uploadSpinner.FinalMSG = fmt.Sprintf("Uploading objects: %d/%d done  %s\n", o.uploaded, o.delta, ui.DoneCheck())
	o.uploadSpinner.Stop()
}

func (o *DeployObserver) listenUploads() {
	for {
		select {
		case <-o.uploadedC:
			o.uploaded += 1
			o.uploadSpinner.Prefix = fmt.Sprintf("Uploading objects: %d/%d ", o.uploaded, o.delta)
		case <-o.uploadedD:
			return
		}
	}
}
