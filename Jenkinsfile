#!/usr/bin/groovy

// load pipeline functions
// Requires pipeline-github-lib plugin to load library from github

@Library('github.com/lachie83/jenkins-pipeline@dev')

def pipeline = new io.estrado.Pipeline()

podTemplate(label: 'jenkins-pipeline', containers: [
    containerTemplate(name: 'jnlp', image: 'lachlanevenson/jnlp-slave:3.10-1-alpine', args: '${computer.jnlpmac} ${computer.name}', workingDir: '/home/jenkins'),
    containerTemplate(name: 'docker', image: 'docker:18.05', command: 'cat', ttyEnabled: true),
    containerTemplate(name: 'golang', image: 'golang:1.10.3', command: 'cat', ttyEnabled: true),
    containerTemplate(name: 'helm', image: 'lachlanevenson/k8s-helm:v2.9.1', command: 'cat', ttyEnabled: true),
    containerTemplate(name: 'kubectl', image: 'lachlanevenson/k8s-kubectl:v1.8.14', command: 'cat', ttyEnabled: true)
],
volumes:[
    hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock'),
]){

  node ('jenkins-pipeline') {

    def pwd = pwd()
    def chart_dir = "${pwd}/charts/s3api"

    checkout scm

    // read in required jenkins workflow config values
    def config = readJSON file: 'Jenkinsfile.json'
    println "pipeline config ==> ${config}"

    // continue only if pipeline enabled
    if (!config.pipeline.enabled) {
        println "pipeline disabled"
        return
    }

    // set additional git envvars for image tagging
    pipeline.gitEnvVars()

    // If pipeline debugging enabled
    if (config.pipeline.debug) {
      println "DEBUG ENABLED"
      sh "env | sort"

      println "Runing kubectl/helm tests"
      container('kubectl') {
        pipeline.kubectlTest()
      }
      container('helm') {
        pipeline.helmConfig()
      }
    }

    def acct = pipeline.getContainerRepoAcct(config)

    // tag image with version, and branch-commit_id
    def image_tags_map = pipeline.getContainerTags(config)

    // compile tag list
    def image_tags_list = pipeline.getMapValues(image_tags_map)

    container('golang') {
        stage ('prepare') {
            // Move source to GOPATH
            sh """
            mkdir -p \$GOPATH/src/github.com/sythe21
            ln -s \$(realpath .) \$GOPATH/src/github.com/sythe21
            """
        }
        stage ('download dependencies') {
            sh """
            cd \${GOPATH}/src/github.com/sythe21/s3api
            go get -u github.com/golang/dep/cmd/dep
            dep ensure
            """
          }

        stage ('compile and test') {
            sh """
            cd \${GOPATH}/src/github.com/sythe21/s3api
            go test -v -race ./...
            """
        }
    }


    stage ('test deployment') {

      container('helm') {

        // run helm chart linter
        pipeline.helmLint(chart_dir)

        // run dry-run helm chart installation
        pipeline.helmDeploy(
          dry_run       : true,
          name          : config.app.name,
          namespace     : config.app.name,
          chart_dir     : chart_dir,
          set           : [
            "imageTag": image_tags_list.get(0),
            "replicas": config.app.replicas,
            "cpu": config.app.cpu,
            "memory": config.app.memory,
            "ingress.hostname": config.app.hostname,
          ]
        )

      }
    }

    stage ('publish container') {

      container('docker') {

        // perform docker login to container registry as the docker-pipeline-plugin doesn't work with the next auth json format
        withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: config.container_repo.jenkins_creds_id, usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
          sh "docker login -u ${env.USERNAME} -p ${env.PASSWORD} ${config.container_repo.host}"
        }

        println "Running Docker build/publish: ${config.container_repo.host}/${config.container_repo.acct}/${config.container_repo.repo}:${config.container_repo.tags}"

        docker.withRegistry("${config.container_repo.host}", "${config.container_repo..auth_id}") {

            def img = docker.image("${config.container_repo.acct}/${config.container_repo.repo}")
            sh "docker build --build-arg VCS_REF=${env.GIT_SHA} --build-arg BUILD_DATE=`date -u +'%Y-%m-%dT%H:%M:%SZ'` -t ${config.container_repo.acct}/${config.container_repo.repo} ."
            for (int i = 0; i < config.container_repo.tags.size(); i++) {
                img.push(config.container_repo.tags.get(i))
            }

            return img.id
        }

        // build and publish container
        pipeline.containerBuildPub(
            dockerfile: config.container_repo.dockerfile,
            host      : config.container_repo.host,
            acct      : acct,
            repo      : config.container_repo.repo,
            tags      : image_tags_list,
            auth_id   : config.container_repo.jenkins_creds_id,
            image_scanning: config.container_repo.image_scanning
        )

        // anchore image scanning configuration
        println "Add container image tags to anchore scanning list"
        
        def tag = image_tags_list.get(0)
        def imageLine = "${config.container_repo.host}/${acct}/${config.container_repo.repo}:${tag}" + ' ' + env.WORKSPACE + '/Dockerfile'
        writeFile file: 'anchore_images', text: imageLine
        anchore name: 'anchore_images', inputQueries: [[query: 'list-packages all'], [query: 'list-files all'], [query: 'cve-scan all'], [query: 'show-pkg-diffs base']]

      }

    }

    if (env.BRANCH_NAME =~ "PR-*" ) {
      stage ('deploy to k8s') {
        container('helm') {
          // Deploy using Helm chart
          pipeline.helmDeploy(
            dry_run       : false,
            name          : env.BRANCH_NAME.toLowerCase(),
            namespace     : env.BRANCH_NAME.toLowerCase(),
            chart_dir     : chart_dir,
            set           : [
              "imageTag": image_tags_list.get(0),
              "replicas": config.app.replicas,
              "cpu": config.app.cpu,
              "memory": config.app.memory,
              "ingress.hostname": config.app.hostname,
            ]
          )

          //  Run helm tests
          if (config.app.test) {
            pipeline.helmTest(
              name        : env.BRANCH_NAME.toLowerCase()
            )
          }

          // delete test deployment
          pipeline.helmDelete(
              name       : env.BRANCH_NAME.toLowerCase()
          )
        }
      }
    }

    // deploy only the master branch
    if (env.BRANCH_NAME == 'master') {
      stage ('deploy to k8s') {
        container('helm') {
          // Deploy using Helm chart
          pipeline.helmDeploy(
            dry_run       : false,
            name          : config.app.name,
            namespace     : config.app.name,
            chart_dir     : chart_dir,
            set           : [
              "imageTag": image_tags_list.get(0),
              "replicas": config.app.replicas,
              "cpu": config.app.cpu,
              "memory": config.app.memory,
              "ingress.hostname": config.app.hostname,
            ]
          )
          
          //  Run helm tests
          if (config.app.test) {
            pipeline.helmTest(
              name          : config.app.name
            )
          }
        }
      }
    }
  }
}
