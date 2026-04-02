package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func createJar(path string) {
	jarUrl := "https://fill-data.papermc.io/v1/objects/158703f75a26f842ea656b3dc6d75bf3d1ec176b97a2c36384d0b80b3871af53/paper-1.21.10-130.jar"
	resp, err := http.Get(jarUrl)

	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(path, data, 0o777)
}
