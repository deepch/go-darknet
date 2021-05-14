package darknet

// #include <darknet.h>
import "C"
import (
	"errors"
	"image"
	"unsafe"
	//"log"
)

// Image represents the image buffer.
type Image struct {
	Width  int
	Height int

	image C.image
}

var errUnableToLoadImage = errors.New("unable to load image")

// Close and release resources.
func (img *Image) Close() error {
	// C.free_src_data(*img)
	//C.free_src_data(img.image)
	//img.image.data = 
//	img.image = nil
	C.free_image(img.image)
	//C.free_image(*img.image)
	return nil
}

// ImageFromPath reads image file specified by path.

func ImageFromPath(path string) (*Image, error) {
	p := C.CString(path)
	defer C.free(unsafe.Pointer(p))

	img := Image{
		image: C.load_image_color(p, 0, 0),
	}

	if img.image.data == nil {
		return nil, errUnableToLoadImage
	}

	img.Width = int(img.image.w)
	img.Height = int(img.image.h)

	return &img, nil
}

func ImageFromMemory(buf []byte) (*Image, error) {
	cBuf := C.CBytes(buf)
	defer C.free(cBuf)

	img := Image{
		image: C.load_image_from_memory_color((*C.char)(cBuf),
			C.int(len(buf)), 0, 0),
	}

	if img.image.data == nil {
		return nil, errUnableToLoadImage
	}

	img.Width = int(img.image.w)
	img.Height = int(img.image.h)

	return &img, nil
}


func DrawYCbCrNew(src *image.YCbCr) *float32 {
	w := src.Bounds().Dx()
	h := src.Bounds().Dy()
	buffer := make([]float32, w*h*3)
	x1 := w * 4
	var all int
	for sy := 0; sy != h; sy++ {
		yi := (sy-src.Rect.Min.Y)*src.YStride + (0 - src.Rect.Min.X)
		ciBase := (sy/2-src.Rect.Min.Y/2)*src.CStride - src.Rect.Min.X/2
		for x, sx := 0, 0; x != x1; x, sx, yi = x+4, sx+1, yi+1 {
			ci := ciBase + sx/2
			yy1 := int32(src.Y[yi]) * 0x10101
			cb1 := int32(src.Cb[ci]) - 128
			cr1 := int32(src.Cr[ci]) - 128
			r := yy1 + 91881*cr1
			if uint32(r)&0xff000000 == 0 {
				r >>= 16
			} else {
				r = ^(r >> 31)
			}
			g := yy1 - 22554*cb1 - 46802*cr1
			if uint32(g)&0xff000000 == 0 {
				g >>= 16
			} else {
				g = ^(g >> 31)
			}
			b := yy1 + 116130*cb1
			if uint32(b)&0xff000000 == 0 {
				b >>= 16
			} else {
				b = ^(b >> 31)
			}
			buffer[all] = float32(r) / 255.
			buffer[w*h+all] = float32(g) / 255.
			buffer[w*h*2+all] = float32(b) / 255.
			all++
		}
	}
	return &buffer[0]
}
func DrawYCbCr(src *image.YCbCr) []C.float {
	w := src.Bounds().Dx()
	h := src.Bounds().Dy()
	buffer := make([]C.float, w*h*3)
	x1 := w * 4
	var all int
	for sy := 0; sy != h; sy++ {
		yi := (sy-src.Rect.Min.Y)*src.YStride + (0 - src.Rect.Min.X)
		ciBase := (sy/2-src.Rect.Min.Y/2)*src.CStride - src.Rect.Min.X/2
		for x, sx := 0, 0; x != x1; x, sx, yi = x+4, sx+1, yi+1 {
			ci := ciBase + sx/2
			yy1 := int32(src.Y[yi]) * 0x10101
			cb1 := int32(src.Cb[ci]) - 128
			cr1 := int32(src.Cr[ci]) - 128
			r := yy1 + 91881*cr1
			if uint32(r)&0xff000000 == 0 {
				r >>= 16
			} else {
				r = ^(r >> 31)
			}
			g := yy1 - 22554*cb1 - 46802*cr1
			if uint32(g)&0xff000000 == 0 {
				g >>= 16
			} else {
				g = ^(g >> 31)
			}
			b := yy1 + 116130*cb1
			if uint32(b)&0xff000000 == 0 {
				b >>= 16
			} else {
				b = ^(b >> 31)
			}
			buffer[all] = C.float(float32(r) / 255.)
			buffer[w*h+all] = C.float(float32(g) / 255.)
			buffer[w*h*2+all] = C.float(float32(b) / 255.)
			all++
		}
	}
	return buffer
}
/*
func ImageFromFloat3(Y []byte, CB []byte, CR []byte, x, y int) (*Image, error) {
    cY := C.CBytes(Y)
    defer C.free(cY)
    cCB := C.CBytes(CB)
    defer C.free(cCB)
    cCR := C.CBytes(CR)
    defer C.free(cCR)
    img := Image{	
	image: C.load_image_from_memory_color_v2((*C.char)(cY), C.int(len(Y)), (*C.char)(cCB), C.int(len(CB)), (*C.char)(cCR), C.int(len(CR)), C.int(x), C.int(y)),
    }
    if img.image.data == nil {
	return nil, errUnableToLoadImage
    }
    img.Width = int(img.image.w)
    img.Height = int(img.image.h)
    return &img, nil
}
*/
//image load_image_from_memory_stb_v6(char *puc_y, char *puc_u, char *puc_v, int width_y, int height_y)
func ImageFromFloat6(Y []byte, CB []byte, CR []byte, x, y int) (*Image, error) {
    cY := C.CBytes(Y)
    defer C.free(cY)
    cCB := C.CBytes(CB)
    defer C.free(cCB)
    cCR := C.CBytes(CR)
    defer C.free(cCR)
    img := Image{
	image: C.load_image_from_memory_stb_v6((*C.char)(cY), (*C.char)(cCB), (*C.char)(cCR), C.int(x), C.int(y)),
    }
    if img.image.data == nil {
	return nil, errUnableToLoadImage
    }
    
    img.Width = int(img.image.w)
    img.Height = int(img.image.h)
//    log.Println(img.Width, img.Height)
    return &img, nil
}
/*
func ImageFromFloat7(y, u, v []byte, w, h int) (*Image, error) {
    yc := C.CBytes(y)
    uc := C.CBytes(u)
    vc := C.CBytes(v)
    defer C.free(yc)
    defer C.free(uc)
    defer C.free(vc)
    //imgr := C.image{}
    //imgr = C.make_image(C.int(x), C.int(y), C.int(3));
    //log.Println((*imgr.data)[1])
    img := Image{
	image: C.load_image_from_memory_stb_v7((*C.char)(yc),(*C.char)(uc),(*C.char)(vc), C.int(w), C.int(h)),
//	image: imgr,
    }
    if img.image.data == nil {
	return nil, errUnableToLoadImage
    }
    img.Width = int(img.image.w)
    img.Height = int(img.image.h)
    return &img, nil
}
*/
func ImageFromFloat3(Y []byte, CB []byte, CR []byte, x, y int) (*Image, error) {
    cY := C.CBytes(Y)
    defer C.free(cY)
    cCB := C.CBytes(CB)
    defer C.free(cCB)
    cCR := C.CBytes(CR)
    defer C.free(cCR)
    img := Image{
	image: C.load_image_from_memory_color_v2((*C.char)(cY), C.int(len(Y)), (*C.char)(cCB), C.int(len(CB)), (*C.char)(cCR), C.int(len(CR)), C.int(x), C.int(y)),
    }
    if img.image.data == nil {
	return nil, errUnableToLoadImage
    }
    img.Width = int(img.image.w)
    img.Height = int(img.image.h)
    return &img, nil
}
func ImageFromFloat5(buf []byte, w, h int) (*Image, error) {
    cBuf := C.CBytes(buf)
    //(*C.char)(cBuf),
    //<------><------><------>C.int(len(buf))
    img := Image{
	image: C.load_image_from_memory_color_v4((*C.char)(cBuf), C.int(len(buf)), C.int(w), C.int(h)),
    }
    if img.image.data == nil {
	return nil, errUnableToLoadImage
    }
    img.Width = int(img.image.w)
    img.Height = int(img.image.h)
    return &img, nil
}
func ImageFromFloat4(buf []float32, w, h, rw, rh int) (*Image, error) {
    img := Image{
	image: C.load_image_from_memory_color_v3((*C.float)(unsafe.Pointer(&buf[0])), C.int(w), C.int(h)),
    }
    if img.image.data == nil {
	return nil, errUnableToLoadImage
    }
    img.Width = int(img.image.w)
    img.Height = int(img.image.h)
    return &img, nil
}
/*
func ImageFromFloat2(buf *float32, w, h int) (*Image, error) {
	//	cBuf := C.CBytes(buf)
	//defer C.free(cBuf)
	//	ddd := time.Now()
	//buf := DrawYCbCrNew(imgs)
	//log.Println("old", time.Now().Sub(ddd))
	//ddd = time.Now()
	//DrawYCbCrNew(imgs)
	//log.Println("new", time.Now().Sub(ddd))
	data := C.alloc_src_data()

	data.h = C.int(h)
	data.w = C.int(w)
	data.c = C.int(3)
	//data.data = &(*buf)[0]

	data.data = (*C.float)(unsafe.Pointer(buf))
	img := Image{
		image: data,
	}
	//img.image.data = &buf[0]
	//in := []C.float{1.23, 4.56}
	//	img := Image{C.make_image_deep(C.int(w), C.int(h), (*C.float)(&buf[0]))}
	if img.image.data == nil {
		return nil, errUnableToLoadImage
	}
	//C.free_src_data(data)
	img.Width = int(img.image.w)
	img.Height = int(img.image.h)
	//	log.Fatalln("тут всё ок")
	return &img, nil
}
*/
/*
func ImageFromFloat(buf []byte, w, h int) (*Image, error) {
	cBuf := C.CBytes(buf)
	//defer C.free(cBuf)

	img := Image{
		image: C.load_image_from_memory_deep((*C.uchar)(cBuf),
			C.int(len(buf)), C.int(w), C.int(h), 4),
	}

	if img.image.data == nil {
		return nil, errUnableToLoadImage
	}

	img.Width = int(img.image.w)
	img.Height = int(img.image.h)
	return &img, nil
}
*/
/*
func ImageFromFloat(buf []byte, w, h int) (*Image, error) {
	cBuf := C.CBytes(buf)
	//	defer C.free(cBuf)
	//a := buf
	//	in := []C.float{1.23, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56, 4.56}
	//t := (*C.float)(C.calloc(C.size_t(2*2*4), 6))
	//defer C.free(unsafe.Pointer(t))
	//var a [16]float32
	//	C.getMatrix((**C.float)(unsafe.Pointer(&a[0])))
	img := Image{
		//float_to_image(int w, int h, int c, float *data)
		image: C.make_image(C.int(w), C.int(h), 3),
		//image: C.make_random_image(100, 100, 4),
	}
	img.image.w = C.int(w)
	img.image.h = C.int(h)
	img.image.c = 4
	//	img.image.data = (*C.float)(cBuf)
	C.write_byte_buffer_as_ratio_to_image_data((*C.char)(cBuf), C.int(w), C.int(h), 3)
	if img.image.data == nil {
		return nil, errUnableToLoadImage
	}
	//	C.yuv_to_rgb(img.image)
	img.Width = int(w)
	img.Height = int(h)

	return &img, nil
}
*/
