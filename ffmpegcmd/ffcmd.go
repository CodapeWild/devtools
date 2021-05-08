package ffmpegcmd

import (
	"devtools/clock"
	"devtools/comerr"
	"devtools/command"
	"devtools/file"
	"fmt"
	"log"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type FFStream struct {
	Index         int    `json:"index"`
	CodecName     string `json:"codec_name"`
	CodecLongName string `json:"codec_long_name"`
	CodecType     string `json:"codec_type"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
}

type FFInfo struct {
	Streams []*FFStream `json:"streams"`
	Format  struct {
		FileName       string `json:"filename"`
		NbStreams      int    `json:"nb_streams"`
		NbPrograms     int    `json:"nb_programs"`
		FormatName     string `json:"format_name"`
		FormatLongName string `json:"format_long_name"`
		StartTime      string `json:"start_time"`
		Duration       string `json:"duration"`
		Span           int64
		Size           string `json:"size"`
		BitRate        string `json:"bit_rate"`
	} `json:"format"`
}

func ParseFFInfoFromFile(filePath string) (*FFInfo, error) {
	if !file.IsFileExists(filePath) {
		return nil, comerr.ErrParamInvalid
	}

	ffinfo := &FFInfo{}
	cmdstr := "ffprobe -v quiet -print_format json -show_format -show_streams " + filePath
	cmdary := strings.Split(cmdstr, " ")
	cmd := &(command.Command{Cmd: exec.Command(cmdary[0], cmdary[1:]...)})
	if err := cmd.StdoutJson(ffinfo); err != nil {
		return nil, err
	} else {
		if i := strings.Index(ffinfo.Format.Duration, "."); i > 0 {
			ffinfo.Format.Span, err = strconv.ParseInt(ffinfo.Format.Duration[:i], 10, 64)
		} else {
			ffinfo.Format.Span, err = strconv.ParseInt(ffinfo.Format.Duration, 10, 64)
		}

		return ffinfo, err
	}
}

type FFCmd struct {
	finfo                     *FFInfo
	ss                        string
	input, output             string
	ivcodec, iacodec          string
	oformat, ovcodec, oacodec string
	iwidth, iheight           int
	owidth, oheight           int
	err                       error
}

func NewFFCmd(finfo *FFInfo) *FFCmd {
	if finfo == nil {
		return nil
	}

	ffcmd := &FFCmd{finfo: finfo}
	ffcmd.input = finfo.Format.FileName
	for _, v := range ffcmd.finfo.Streams {
		switch v.CodecType {
		case "video":
			ffcmd.ivcodec = v.CodecName
			ffcmd.iwidth = v.Width
			ffcmd.iheight = v.Height
		case "audio":
			ffcmd.iacodec = v.CodecName
		}
	}

	return ffcmd
}

func (this *FFCmd) AddOutput(filePath string) *FFCmd {
	this.output = filePath
	this.oformat = path.Ext(filePath)[1:]
	if this.oformat == "" {
		this.err = comerr.ErrParamInvalid
	}

	return this
}

func (this *FFCmd) SetTranscode(vcodec, acodec string) *FFCmd {
	this.ovcodec = "copy"
	this.oacodec = "copy"

	if this.ivcodec != vcodec {
		this.ovcodec = vcodec
	}
	if this.iacodec != acodec {
		this.oacodec = acodec
	}

	return this
}

func (this *FFCmd) ScaleTo(width, height int) *FFCmd {
	this.owidth = -1
	this.oheight = -1

	var changed bool
	if width > 0 && width < this.iwidth {
		this.owidth = width
		changed = true
	}
	if height > 0 && height < this.iheight {
		this.oheight = height
		changed = true
	}
	if changed && (this.ovcodec == "" || this.ovcodec == "copy") {
		this.ovcodec = this.ivcodec
	}

	return this
}

// filePath need extension be set
func (this *FFCmd) SetStartOffset(escape int) *FFCmd {
	if int64(escape) > this.finfo.Format.Span {
		this.ss = "00:00:00"
	} else {
		this.ss = clock.FormatSeconds(escape)
	}

	return this
}

// for extract frame the args: -vframes 1 -q:v 3 (for jpeg format, scale of 1-31)
func (this *FFCmd) Extract(args ...string) (*exec.Cmd, error) {
	if this.err != nil {
		return nil, this.err
	}

	cmdstr := fmt.Sprintf("ffmpeg -ss %s -i %s", this.ss, this.input)
	if this.owidth != -1 || this.oheight != -1 {
		cmdstr = fmt.Sprintf("%s -vf scale=%d:%d", cmdstr, this.owidth, this.oheight)
	}
	if len(args) != 0 {
		cmdstr = fmt.Sprintf("%s %s %s", cmdstr, strings.Join(args, " "), this.output)
	} else {
		cmdstr = fmt.Sprintf("%s %s", cmdstr, this.output)
	}

	cmdary := strings.Split(cmdstr, " ")
	log.Println(cmdstr)

	return exec.Command(cmdary[0], cmdary[1:]...), nil
}

// for transcode the args: -crf 18 -preset slow/fast ...
func (this *FFCmd) Transcode(args ...string) (*exec.Cmd, error) {
	if this.err != nil {
		return nil, this.err
	}

	cmdstr := fmt.Sprintf("ffmpeg -i %s", this.input)
	if this.owidth != -1 || this.oheight != -1 {
		cmdstr = fmt.Sprintf("%s -vf scale=%d:%d", cmdstr, this.owidth, this.oheight)
	}

	if this.ovcodec != "" {
		cmdstr = fmt.Sprintf("%s -c:v %s", cmdstr, this.ovcodec)
	} else {
		cmdstr += " -c:v copy"
	}
	if this.oacodec != "" {
		cmdstr = fmt.Sprintf("%s -c:a %s", cmdstr, this.oacodec)
	} else {
		cmdstr += " -c:a copy"
	}
	if len(args) != 0 {
		cmdstr = fmt.Sprintf("%s %s -f %s %s", cmdstr, strings.Join(args, " "), this.oformat, this.output)
	} else {
		cmdstr = fmt.Sprintf("%s -f %s %s", cmdstr, this.oformat, this.output)
	}
	if this.oformat == "mp4" {
		cmdstr += " -movflags +faststart"
	}

	cmdary := strings.Split(cmdstr, " ")
	log.Println(cmdstr)

	return exec.Command(cmdary[0], cmdary[1:]...), nil
}
