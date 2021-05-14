package darknet

// #include <stdlib.h>
//
// #include <darknet.h>
//
// #include "network.h"
import "C"
import (
	"errors"
	"time"
	"unsafe"
	//"log"
)

// YOLONetwork represents a neural network using YOLO.
type YOLONetwork struct {
	GPUDeviceIndex           int
	DataConfigurationFile    string
	NetworkConfigurationFile string
	WeightsFile              string
	Threshold                float32

	ClassNames []string
	Classes    int

	cNet                *C.network
	hierarchalThreshold float32
	nms                 float32
}

var errNetworkNotInit = errors.New("network not initialised")
var errUnableToInitNetwork = errors.New("unable to initialise")

// Init the network.
func (n *YOLONetwork) Init() error {
	nCfg := C.CString(n.NetworkConfigurationFile)
	defer C.free(unsafe.Pointer(nCfg))
	wFile := C.CString(n.WeightsFile)
	defer C.free(unsafe.Pointer(wFile))

	// GPU device ID must be set before `load_network()` is invoked.
	C.cuda_set_device(C.int(n.GPUDeviceIndex))
	n.cNet = C.load_network(nCfg, wFile, 0)

	if n.cNet == nil {
		return errUnableToInitNetwork
	}

	C.set_batch_network(n.cNet, 1)
	C.srand(2222222)

	// Currently, hierarchal thresholdcuda is always 0.5.
	n.hierarchalThreshold = .5

	// Currently NMS is always 0.45.
	n.nms = .45

	n.Classes = int(C.get_network_layer_classes(n.cNet, n.cNet.n-1))
	cClassNames := loadClassNames(n.DataConfigurationFile)
	defer freeClassNames(cClassNames)
	n.ClassNames = makeClassNames(cClassNames, n.Classes)

	return nil
}

// Close and release resources.
func (n *YOLONetwork) Close() error {
	if n.cNet == nil {
		return errNetworkNotInit
	}

	C.free_network(n.cNet)
	n.cNet = nil
	return nil
}

// Detect specified image.
func (n *YOLONetwork) Detect(img *Image) (*DetectionResult, error) {
	if n.cNet == nil {
		return nil, errNetworkNotInit
	}

	startTime := time.Now()
	//log.Println("1", n, n.cNet, n.Classes)
	//log.Println(n.Classes, n.Threshold, n.hierarchalThreshold)
	// n.Threshold = 0.1
	//n.hierarchalThreshold = .1
	result := C.perform_network_detect(n.cNet, &img.image, C.int(n.Classes),
		C.float(n.Threshold), C.float(n.hierarchalThreshold), C.float(n.nms))
	//log.Println(int(result.detections_len))
	endTime := time.Now()
	defer C.free_detections(result.detections, result.detections_len)
	//log.Println(result)
	ds := makeDetections(img, result.detections, int(result.detections_len),
		n.Threshold, n.Classes, n.ClassNames)
	//log.Println("3")
	endTimeOverall := time.Now()

	out := DetectionResult{
		Detections:           ds,
		NetworkOnlyTimeTaken: endTime.Sub(startTime),
		OverallTimeTaken:     endTimeOverall.Sub(startTime),
	}

	return &out, nil
}
