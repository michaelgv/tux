package main

type ConfigurationFile struct {
	Useredis string `yaml:"useredis"`
	STATIC_FOLDER string `yaml:"static_folder"`
}