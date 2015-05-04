package image

import (
	"github.com/nfnt/resize"

	"gopkg.in/mgo.v2/bson"

	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"

	"github.com/PrincetonOBO/OBOBackend/util"
)

const FULL string = "full"
const THUMB string = "thumb"
const WIDTH uint = 100
const HEIGHT uint = 100

type Image struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Item_Id   bson.ObjectId `json:"item_id"`
	Image     string        `json:"image"`
	Thumbnail string        `json:"thumbnail"`
}

func NewImage(raw image.Image, itemId bson.ObjectId) *Image {
	fullJpeg := new(bytes.Buffer)
	jpeg.Encode(fullJpeg, raw, nil)
	imString := base64.StdEncoding.EncodeToString(fullJpeg.Bytes())

	scaledIm := resize.Thumbnail(WIDTH, HEIGHT, raw, resize.MitchellNetravali)
	scaledJpeg := new(bytes.Buffer)
	jpeg.Encode(scaledJpeg, scaledIm, nil)
	thumbString := base64.StdEncoding.EncodeToString(scaledJpeg.Bytes())

	im := new(Image)
	im.Item_Id = itemId
	im.Image = imString
	im.Thumbnail = thumbString

	return im

}

type ImagePresenter struct {
	Item_Id bson.ObjectId `json:"item_id"`
	Image   string        `json:"image"`
	Size    string        `json:"size"`
}

func (i Image) ToPresenter(size string) ImagePresenter {
	imPres := ImagePresenter{Item_Id: i.Item_Id}
	if size == FULL {
		imPres.Size = FULL
		imPres.Image = i.Image

	} else {
		imPres.Size = THUMB
		imPres.Image = i.Thumbnail
	}
	return imPres
}

func (i *ImagePresenter) ToImage() Image {
	var imBytes []byte
	//util.Log(i.Image)
	imBytes, err1 := base64.StdEncoding.DecodeString(i.Image)
	util.Logerr(err1)
	buf := bytes.NewBuffer(imBytes)
	//util.Log(imBytes)
	raw, err := jpeg.Decode(buf)
	util.Logerr(err)
	im := NewImage(raw, i.Item_Id)
	return *im
}
