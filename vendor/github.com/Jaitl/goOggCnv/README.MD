# goOggCnv
Converts mp3 and wav in ogg opus format. Wrapper around *ffmpeg*.

## Use Case
Converting audio for voice messages in Telegram.

## Usage
create with explicit path to *ffmpeg*:
```
cnv := goOggCnv.New("path to ffmpeg")
```
or create with default system path to *ffmpeg*:
```
cnv := goOggCnv.NewD()
```
convert:
```
err := cnv.Mp3ToOgg(inputMp3Path, outputOpusPath)
```
```
err := cnv.WavToOgg(inputWavPath, outputOpusPath)
```

## Install
1. Install `ffmpeg`
2. Install *goOggCnv* library:
```
go get github.com/Jaitl/goOggCnv
```
