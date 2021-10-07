#!/usr/bin/env groovy

String image_name = 'simple-go-helloworld'

node {
  docker.image('nimmis/alpine-golang').inside('-u root') {
    // Preparing container
    stage ('System requirenments') {
      checkout scm
      
      // installing system required packages
      sh '''
        apk add --update libltdl git make
      '''
      // symlink to GOPATH and move to application workspace
      sh '''
        mkdir -p $GOPATH/src
        ln -s $WORKSPACE $GOPATH/src
        cd $GOPATH/src/simple-go-helloworld
      '''
    }
    
    // Run code testing
    stage('Test') {
      sh '''
        make test
      '''
    }
    
    // build binary
    stage('Build') {
      // build the code
      sh '''
        make build
      '''
    }
    
    // deploy
    stage('Deploy') {
      // build image
      if (params.DOCKER_REGISTRY) {
        image_name = ["${params.DOCKER_REGISTRY}","${image_name}"].join('/')
      }
        
      // build image
      String short_commit = sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim()
      echo "Building image: ${short_commit}"
      echo sh(returnStdout: true, script: "docker build -t ${image_name}:${short_commit} .").trim()
      echo sh(returnStdout: true, script: "docker tag ${image_name}:${short_commit} ${image_name}:latest").trim()
      if (params.DOCKER_REGISTRY) {
        echo sh(returnStdout: true, script: "docker push ${image_name}:${short_commit}").trim()
        echo sh(returnStdout: true, script: "docker push ${image_name}:latest").trim()
      }
    }
  }  
}