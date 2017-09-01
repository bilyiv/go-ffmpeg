package go_ffmpeg

import (
	"fmt"
	os_exec "os/exec"
)

type Video struct {
	Path string
}

// ExtractAudio extracts audio from video.
func (v Video) ExtractAudio(saveTo string, bps int, ch int, fr int) error {
	template := "ffmpeg -y -i %s -ab %d -ac %d -ar %d -vn %s"
	return exec(template, v.Path, bps, ch, fr, saveTo)
}

// CombineAudio combines video's audio with additional audio then save the result to a new video.
func (v Video) CombineAudio(saveTo string, audioPath string) error {
	template := "ffmpeg -y -i %s -i %s -filter_complex \"[0:a][1:a]amerge[aout]\" -map 0:v -map \"[aout]\" -c:v copy -shortest %s"
	return exec(template, v.Path, audioPath, saveTo)
}

// MakeScreenshot makes a screenshot of a specific second.
func (v Video) MakeScreenshot(saveTo string, second int) error {
	template := "ffmpeg -y -i %s -ss %d -qscale:v 2 -vframes 1 %s"
	return exec(template, v.Path, second, saveTo)
}

// SacelAndCrop scales then crops video to specific size.
func (v Video) ScaleAndCrop(saveTo string, width int, height int, vb int, ab int) error {
	template := "ffmpeg -y -i %s -vf \"scale=(iw*sar)*max(%d/(iw*sar)\\,%d/ih):ih*max(%d/(iw*sar)\\,%d/ih), crop=%d:%d\" -vb %d -ab %d -s %dx%d %s"
	return exec(template, v.Path, width, height, width, height, width, height, vb, ab, width, height, saveTo)
}

// Wrap wraps video by other videos.
func (v Video) Wrap(saveTo string, begin string, end string) error {
	template := "ffmpeg -y -i %s -i %s -i %s -filter_complex \"[0:v:0] [0:a:0] [1:v:0] [1:a:0] [2:v:0] [2:a:0] concat=n=3:v=1:a=1 [v] [a]\" -map \"[v]\" -map \"[a]\" %s"
	return exec(template, begin, v.Path, end, saveTo)
}

// Exec executes command.
func exec(template string, params ...interface{}) error {
	cmd := fmt.Sprintf(template, params...)
	return os_exec.Command("sh", "-c", cmd).Run()
}
