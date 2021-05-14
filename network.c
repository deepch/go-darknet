#include <stdlib.h>

#include <darknet.h>

#include "network.h"

int get_network_layer_classes(network *n, int index)
{
	return n->layers[index].classes;
}

struct network_box_result perform_network_detect(network *n, image *img,
    int classes, float thresh, float hier_thresh, float nms)
{

//char *filename = "deep.png";
//char buff[256];
//    char *input = buff;

 //if(filename){
//            strncpy(input, filename, 256);
//        } else {
//            printf("Enter Image Path: ");
//            fflush(stdout);
//            input = fgets(input, 256, stdin);
//            if(!input) return;
//            strtok(input, "\n");
//        }
//        image im = load_image_color(input,0,0);


    image sized = letterbox_image(*img, n->w, n->h);
//image sized = letterbox_image(im, n->w, n->h);
    struct network_box_result result = { NULL };

    //  float *X = img->data;
     float *X = sized.data;
//    network_predict(n, X);
    network_predict(n, X);
    result.detections = get_network_boxes(n, img->w, img->h, thresh, hier_thresh, 0, 1, &result.detections_len);
    //result.detections = get_network_boxes(n, 416, 416, thresh, hier_thresh, 0, 1, &result.detections_len); 
    //printf("found %d elements\n", result.detections_len);
    //printf("ffff= %.6f \r\n", result.detections_len);
    if (nms) {
        do_nms_sort(result.detections, result.detections_len, classes, nms);
    }
    free_image(sized);
    //printf("ffff= %u %u %.6f %.6f %.6f %u \r\n", img->w, img->h, thresh, hier_thresh, nms, result.detections_len);
    return result;
}
