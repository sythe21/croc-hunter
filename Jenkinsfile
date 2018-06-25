#!/usr/bin/groovy

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
    def go_dir = "github.com/sythe21/s3api"

    git 'https://github.com/sythe21/s3api.git'

    // read in required jenkins workflow config values
    def config = readYaml file: 'Jenkinsfile.yml'
    println "pipeline config ==> ${config}"

    container('golang') {
        stage ('prepare') {
            // Move source to GOPATH
            sh """
            mkdir -p \$GOPATH/src/github.com/sythe21
            ln -s \$(realpath .) \$GOPATH/src/github.com/sythe21
            """
        }
        stage ('go dependencies') {
            sh "cd \$GOPATH/src/${go_dir} && make init-build && make dep"
          }

        stage ('go build') {
            sh "cd \$GOPATH/src/${go_dir} && make build"
        }

        stage ('go test') {
            sh "cd \$GOPATH/src/${go_dir} && make test"
        }
    }

    stage ('helm verify') {

      container('helm') {

        sh """
        helm lint ${chart_dir}
        helm upgrade --dry-run --install ${config.name} ${chart_dir} --namespace=default --values jenkins-deploy.yml
        """
      }
    }

    container('docker') {
        stage ('docker login') {
            // perform docker login to container registry as the docker-pipeline-plugin doesn't work with the next auth json format
            withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: config.image.jenkinsCredsId, usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS']]) {
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
            println "Pushing docker images to repository ${config.image.name}"
            sh "make push"
        }
    }

    stage ('deploy to kubernetes') {
        container('helm') {
            sh """
            println "Running deployment"
            sh "helm upgrade --install ${config.name} ${chart_dir} --namespace=default --values jenkins-deploy.yml --wait"
            println "Application ${config.name} successfully deployed"
            """
        }
    }
  }
}
