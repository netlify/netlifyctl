package operations

import (
	"fmt"

	"github.com/netlify/netlifyctl/ui"
	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/porcelain"
)

const progressColor = "blue"

type DeployObserver struct {
	walkerSpinner *ui.TaskTracker
	deltaSpinner  *ui.TaskTracker
	uploadSpinner *ui.TaskTracker

	total    int
	delta    int
	uploaded int

	uploadedC chan *porcelain.FileBundle
	uploadedD chan bool
}

func NewDeployObserver() *DeployObserver {
	return &DeployObserver{
		walkerSpinner: ui.NewTaskTracker(),
		deltaSpinner:  ui.NewTaskTracker(),
		uploadSpinner: ui.NewTaskTracker(),
		uploadedC:     make(chan *porcelain.FileBundle),
		uploadedD:     make(chan bool),
	}
}

func (o *DeployObserver) OnSetupWalk() error {
	o.walkerSpinner.Start("Counting objects .... ")
	return nil
}

func (o *DeployObserver) OnSuccessfulStep(*porcelain.FileBundle) error {
	o.total += 1
	o.walkerSpinner.Step(fmt.Sprintf("Counting objects: %d ", o.total))
	return nil
}

func (o *DeployObserver) OnSuccessfulWalk(df *models.DeployFiles) error {
	o.walkerSpinner.Success(fmt.Sprintf("Counting objects: %d total objects", o.total))
	return nil
}

func (o *DeployObserver) OnFailedWalk() {
	o.walkerSpinner.Failure("Counting objects")
}

func (o *DeployObserver) OnSetupDelta(*models.DeployFiles) error {
	o.deltaSpinner.Start("Resolving deltas .... ")
	return nil
}

func (o *DeployObserver) OnSuccessfulDelta(df *models.DeployFiles, d *models.Deploy) error {
	o.delta = len(d.Required) + len(d.RequiredFunctions)

	msg := fmt.Sprintf("Resolving deltas: %d objects to upload", o.delta)
	o.deltaSpinner.Success(msg)

	go o.listenUploads()
	msg = fmt.Sprintf("Uploading objects: %d/%d ", o.uploaded, o.delta)
	o.uploadSpinner.Start(msg)
	return nil
}

func (o *DeployObserver) OnFailedDelta(*models.DeployFiles) {
	o.deltaSpinner.Failure("Resolving deltas")
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
	o.uploadSpinner.Failure("Uploading objects")
}

func (o *DeployObserver) Finish() {
	o.uploadedD <- true
	o.uploadSpinner.Success(fmt.Sprintf("Uploading objects: %d/%d done", o.uploaded, o.delta))
}

func (o *DeployObserver) listenUploads() {
	for {
		select {
		case <-o.uploadedC:
			o.uploaded += 1
			o.uploadSpinner.Step(fmt.Sprintf("Uploading objects: %d/%d ", o.uploaded, o.delta))
		case <-o.uploadedD:
			return
		}
	}
}
