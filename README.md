# `trans_android_str`

[![Build Status](https://travis-ci.org/northbright/trans_android_str.svg?branch=master)](https://travis-ci.org/northbright/trans_android_str)

`trans_android_str` is a tool written in [Golang](http://golang.org) which generates translated string xml files under Android app resource path(xx/res/values-xx) by given translated strings and config file.

#### How it works

* Translation File Format:

    * The translated strings of each language start with the iso 639-1 language name(and iso 3166-1 locale name if need).
    * Each line contains a translated string.
    * A blank line need to be added at the end of last translated string.
    * UTF-8 without BOM

Ex:  
We need to translate 3 strings:  
Music, Movie, Voice
in original string xml:

`res/values/strings.xml`:

    // res/values/strings.xml
    <?xml version="1.0" encoding="utf-8"?>
    <resources>
        <string name="yes">"Yes"</string>
        <string name="no">"No"</string>
        <string name="music">Music</string>
        <add-resource type="string" name="movie" />
        <string name="movie">Movie</string>
        <add-resource type="string" name="voice" />
        <string name="voice">Voice</string>
    </resources>

`translation.txt`:

    fr
    Musique
    Film
    Voix

    de
    Musik
    Film
    Sprache

    zh-rCN
    音乐
    电影
    语音

    zh-rTW
    音樂
    電影
    語音         

* Configure JSON File Format  
  
    * JSON array contains the string names to need be translated in original strings xml.
    * The count of string names in config file MUST equals to the count of translated strings in translation file.

Ex:  
`config.json`:  
   
    [
        "music",
        "movie",
        "voice"
    ]    

#### Usage:  
`trans_android_str -i <default string xml> -o <out resource folder> -c <config file> -t <translation file>`

Ex:  
`trans_android_str -i "res/values/strings.xml" -o "~/app/res" -c "config.json" -t "translation.txt"`

#### License
* [MIT License](LICENSE)
