package advanced

import (
	"image"
	_ "image/jpeg"
	"log"
	"os"
)

func musnt(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// WithoutPeeVee Advanced example 1 without PeeVee
func WithoutPeeVee() {
	reader, err := os.Open("advanced/landscape.jpg")
	defer reader.Close()
	musnt(err)

	img, _, err := image.Decode(reader)
	musnt(err)

	grayscaledImage := image.NewGray(img.Bounds())

	for col := 0; col < img.Bounds().Max.X; col++ {
		for row := 0; row < img.Bounds().Max.Y; row++ {
			grayscaledImage.Set(row, col, img.At(row, col))
		}
	}

	// for col := 0; col < img.Bounds().Max.X; col++ {
	// 	for row := 0; row < img.Bounds().Max.Y; row++ {
	// 		fmt.Println(row, "/", col, img.At(row, col))
	// 		time.Sleep(time.Millisecond * 200)
	// 	}
	// }
}

// WithPeeVee Advanced example 1 with PeeVee
func WithPeeVee() {

}
