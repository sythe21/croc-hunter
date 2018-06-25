#!/usr/bin/groovy

// load pipeline functions
// Requires pipeline-github-lib plugin to load library from github

@Library('github.com/lachie83/jenkins-pipeline@dev')

def pipeline = new io.estrado.Pipeline()

podTemplate(label: 'jenkins-pipeline', containers: [
    containerTemplate(name: 'jnlp', image: 'lachlanevenson/jnlp-slave:3.10-1-alpine', args: '${computer.jnlpmac} ${computer.name}', workingDir: '/home/jenkins', resourceRequestCpu: '50m'),
    containerTemplate(name: 'docker', image: 'docker:18.05', command: 'cat', ttyEnabled: true, resourceRequestCpu: '50m'),
    containerTemplate(name: 'golang', image: 'golang:1.10.3', command: 'cat', ttyEnabled: true, resourceRequestCpu: '50m'),
    containerTemplate(name: 'helm', image: 'lachlanevenson/k8s-helm:v2.9.1', command: 'cat', ttyEnabled: true, resourceRequestCpu: '50m')
],
volumes:[
    hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock'),
]){

  node ('jenkins-pipeline') {

    def pwd = pwd()
    def chart_dir = "${pwd}/charts/s3api"

    git 'https://github.com/sythe21/s3api.git'

    // read in required jenkins workflow config values
    def config = readJSON file: 'Jenkinsfile.json'
    println "pipeline config ==> ${config}"

    container('golang') {
        stage ('prepare') {
            // Move source to GOPATH
            sh """
            apk add make
            mkdir -p \$GOPATH/src/github.com/sythe21
            ln -s \$(realpath .) \$GOPATH/src/github.com/sythe21
            """
        }
        stage ('go dependencies') {
            sh """
            make init-build
            make dep
            """
          }

        stage ('go build') {
            sh "make build"
        }

        stage ('go test') {
            sh "make test"
        }
    }

    stage ('helm verify') {

      container('helm') {

        sh "helm lint ${chart_dir}"

        // run dry-run helm chart installation
        pipeline.helmDeploy(
          dry_run       : true,
          name          : config.app.name,
          namespace     : config.app.name,
          chart_dir     : chart_dir,
          set           : [
            "name": config.app.name,
            "image.name": config.container_repo.image,
            "image.tag": "latest",
            "image.replicaCount": 1,
            "ingress.enable": config.app.ingressEnable,
            "ingress.hostname": config.app.ingressHostname,
          ]
        )

      }
    }

    container('docker') {
        stage ('docker login') {
            // perform docker login to container registry as the docker-pipeline-plugin doesn't work with the next auth json format
            withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: config.container_repo.jenkins_creds_id, usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS']]) {
              sh "make login"
            }
        }

        stage ('docker build') {
            sh """
            println "Building and tagging docker images"
            make build
            make tag
            """
        }

        stage ('docker push') {
            println "Pushing docker images to repository ${config.container_repo.image}"
            sh "make push"
        }
    }

    stage ('deploy to kubernetes') {
        container('helm') {
            pipeline.helmDeploy(
                dry_run       : false,
                name          : config.app.name,
                namespace     : "default",
                chart_dir     : chart_dir,
                set           : [
                    "name": config.app.name,
                    "image.name": config.container_repo.image,
                    "image.tag": "latest",
                    "image.replicaCount": 1,
                    "ingress.enable": config.app.ingressEnable,
                    "ingress.hostname": config.app.ingressHostname,
                ]
            )
        }
    }
  }
}
