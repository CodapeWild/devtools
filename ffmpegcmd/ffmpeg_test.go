package ffmpegcmd

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestInfo(t *testing.T) {
	ffinfo, err := ParseFFInfoFromFile("./sample/s2.wmv")
	if err != nil {
		log.Panicln(err.Error())
	}
	log.Println(ffinfo.Format)
	log.Println(*(ffinfo.Streams[0]))
}

func TestExtract(t *testing.T) {
	fis, err := ioutil.ReadDir("./sample")
	if err != nil {
		log.Panicln(err.Error())
	}
	for _, v := range fis {
		go func(file string) {
			ffinfo, err := ParseFFInfoFromFile("./sample/" + file)
			if err != nil {
				log.Println(err.Error())
			} else if cmd, err := NewFFCmd(ffinfo).AddOutput("./frame/"+file+".jpeg").ScaleTo(0, 300).SetStartOffset(3).Extract("-vframes 1", "-q:v 6"); err != nil {
				log.Println(err.Error())
			} else {
				log.Println(cmd.Run())
			}
		}(v.Name())
	}
	time.Sleep(time.Second)
}

func TestTranscode(t *testing.T) {
	fis, err := ioutil.ReadDir("./sample")
	if err != nil {
		log.Panicln(err.Error())
	}
	for _, v := range fis {
		go func(file string) {
			ffinfo, err := ParseFFInfoFromFile("./sample/" + file)
			if err != nil {
				log.Println(err.Error())
			} else if cmd, err := NewFFCmd(ffinfo).AddOutput("./transcode/"+file+".mp4").SetTranscode("h264", "aac").ScaleTo(0, 360).Transcode("-crf 18", "-preset fast"); err != nil {
				log.Println(err.Error())
			} else {
				log.Println(cmd.Run())
			}
		}(v.Name())
	}
	time.Sleep(time.Second)
}
