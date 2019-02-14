pipeline {
    agent any

    stages {
        stage('prepare build environment') {
            steps {
                sh('''#!/bin/bash -e
                source ./build-env.sh
                create-go-build-env.sh''')
            }
        }
        stage('Install dependencies') {
            steps {
                sh('''#!/bin/bash -e
                source ./build-env.sh
                glide i''')
            }
        }

        stage('Build') {
            steps {
                withEnv(["DOCKER_JENKINS_HOME=${env.DOCKER_JENKINS_MOUNT}"]) {
                    sh('''#!/bin/bash -e
                    source ./build-env.sh
                    go build''')
                }
            }
        }
        stage('Tests') {
            steps {
                sh('''#!/bin/bash -e
                source ./build-env.sh
                go test simple-relmgt simple-relmgt/cmds/draftcmd simple-relmgt/cmds/checkcmd simple-relmgt/cmds/releasecmd simple-relmgt/cmds/statecmd simple-relmgt/cmds/tagcmd''')
            }
        }
        stage('Deploy') {
            when { branch 'master' }
            steps {
                withCredentials([
                usernamePassword(credentialsId: 'github-jenkins-cred', usernameVariable: 'GITHUB_USER', passwordVariable: 'GITHUB_TOKEN')]) {
                    sh('''#!/bin/bash -e
                    source ./build-env.sh
                    publish.sh latest''')
                }
            }
        }
    }

    post {
        success {
            deleteDir()
        }
    }
}
