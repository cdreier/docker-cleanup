package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {

	logFile, err := os.OpenFile("./out.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logFile.Close()

	if err != nil {
		log.Fatal("cannot open logfile", err.Error())
		return
	}
	log.SetOutput(logFile)

	port := flag.String("port", "8080", "port to start the server on")
	interval := flag.String("interval", "12h", "go duration parsable interval the cleanup job runs")
	flag.Parse()

	tickerDuration, err := time.ParseDuration(*interval)
	if err != nil {
		log.Fatal("wrong interval format")
		return
	}

	// run jobs on start
	exitedContainers()
	danglingImages()
	// and start ticker
	go runTicker(tickerDuration)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logs, err := ioutil.ReadFile("./out.log")
		if err != nil {
			w.Write([]byte("cannot open logfile " + err.Error()))
		} else {
			w.Write(logs)
		}
	})

	http.ListenAndServe(":"+*port, nil)

}

func runTicker(d time.Duration) {
	ticker := time.NewTicker(d)

	for {
		tick := <-ticker.C
		log.Println("TICK", tick.String())
		exitedContainers()
		danglingImages()
	}
}

func exitedContainers() {
	log.Println("RUN: exitedContainers")
	cmd := exec.Command("docker", "ps", "-qa", "--filter", "status=exited", "--no-trunc")
	out, err := cmd.Output()

	if err != nil {
		log.Println("cannot execute 'find exited containers'", err.Error())
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		containerID := scanner.Text()
		log.Println("removing container:", containerID)
		rmCmd := exec.Command("docker", "rm", containerID)
		err := rmCmd.Run()
		if err != nil {
			log.Println("failed to remove image", err.Error())
		}
	}
}

func danglingImages() {
	log.Println("RUN: danglingImages")
	cmd := exec.Command("docker", "images", "--filter", "dangling=true", "-q", "--no-trunc")
	out, err := cmd.Output()

	if err != nil {
		log.Println("cannot execute 'find dangling images'", err.Error())
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		imageID := scanner.Text()
		log.Println("removing image:", imageID)
		rmCmd := exec.Command("docker", "rmi", imageID)
		err := rmCmd.Run()
		if err != nil {
			log.Println("failed to remove image", err.Error())
		}
	}

}

func UNUSED_danglingVolumes() {
	log.Println("RUN: danglingVolumes")
	cmd := exec.Command("docker", "volume", "ls", "-qf", "dangling=true")
	out, err := cmd.Output()

	if err != nil {
		log.Println("cannot execute 'find dangling volumes'", err.Error())
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		log.Println(scanner.Text())
	}
}
