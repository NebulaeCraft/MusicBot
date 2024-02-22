package music

import (
	"MusicBot/config"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
)

type StatVolumeResult struct {
	OriginVolumedB     float64
	OriginVolumeLUFS   float64
	LoudnormVolumeLUFS float64
	LoudnormVolumedB   float64
	GainVolumedB       float64
}

type VolumeData struct {
	Index      int
	MeanVolume float64
}

func ExtractVolumeData(log string) (float64, float64, float64, float64, float64) {
	volumeRegex := regexp.MustCompile(`Parsed_volumedetect_(\d+) @ [^\]]*\] mean_volume: ([-\.\d]+) dB`)
	matches := volumeRegex.FindAllStringSubmatch(log, -1)
	var volumeDataList []VolumeData

	for _, match := range matches {
		index, err := strconv.Atoi(match[1])
		if err != nil {
			fmt.Printf("Error converting index to integer: %s\n", err)
			continue
		}

		meanVolume, err := strconv.ParseFloat(match[2], 64)
		if err != nil {
			fmt.Printf("Error converting mean_volume to float: %s\n", err)
			continue
		}

		volumeDataList = append(volumeDataList, VolumeData{
			Index:      index,
			MeanVolume: meanVolume,
		})
	}

	sort.Slice(volumeDataList, func(i, j int) bool {
		return volumeDataList[i].Index < volumeDataList[j].Index
	})

	inputLUFSRegex := regexp.MustCompile(`Input Integrated:\s+([-\.\d]+) LUFS`)
	outputLUFSRegex := regexp.MustCompile(`Output Integrated:\s+([-\.\d]+) LUFS`)

	inputMatches := inputLUFSRegex.FindAllStringSubmatch(log, -1)
	outputMatches := outputLUFSRegex.FindAllStringSubmatch(log, -1)

	inputLUFS, _ := strconv.ParseFloat(inputMatches[0][1], 64)
	outputLUFS, _ := strconv.ParseFloat(outputMatches[0][1], 64)

	return volumeDataList[0].MeanVolume, inputLUFS, volumeDataList[1].MeanVolume, outputLUFS, volumeDataList[2].MeanVolume
}

func StatVolume(src string) (chan StatVolumeResult, chan string) {
	logger := config.Logger
	logger.Info().Msg(">>> Run FFmpeg StatVolume<<<")
	cmd := exec.Command("ffmpeg",
		"-i",
		src,
		"-af",
		fmt.Sprintf("volumedetect,loudnorm=i=-16:print_format=summary,volumedetect,volume=%ddB,volumedetect", PlayStatus.Volume),
		"-f",
		"null",
		"/dev/null",
	)

	statVolumeResult := make(chan StatVolumeResult, 1)
	statVolumeResultLog := make(chan string, 1)

	go func() {
		out, err := cmd.CombinedOutput()
		if err != nil {
			logger.Error().Err(err).Msg(fmt.Sprintf("Failed to run ffmpeg StatVolume, pid: %d", cmd.Process.Pid))
			//logger.Error().Msg(string(out))
		} else {
			logger.Info().Msg(fmt.Sprintf(">>> FFmpeg StatVolume finished, pid: %d <<<", cmd.Process.Pid))
		}
		inputdB, inputLUFS, loudnormdB, loudnormLUFS, gaindB := ExtractVolumeData(string(out))
		statVolumeResultLog <- string(out)
		statVolumeResult <- StatVolumeResult{
			OriginVolumedB:     inputdB,
			OriginVolumeLUFS:   inputLUFS,
			LoudnormVolumeLUFS: loudnormLUFS,
			LoudnormVolumedB:   loudnormdB,
			GainVolumedB:       gaindB,
		}

	}()
	return statVolumeResult, statVolumeResultLog
}
