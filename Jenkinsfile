def releaseStatus = 3 // Not a release branch
def officialReleaseFileFound = 1
def releaseCmdPath = 'tmp/simple-relmgt'

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

        stage('prepare deployment environment') {
            steps {
                sh('''#!/bin/bash
                mkdir -p tmp
                rm -f tmp/simple-relmgt
                curl -L -s -O tmp/simple-relmgt https://github.com/forj-oss/simple-relmgt/releases/download/latest/simple-relmgt >/dev/null
                if [[ -f tmp/simple-relmgt ]]
                then
                    chmod +x tmp/simple-relmgt
                    tmp/simple-relmgt --version
                else
                    echo "No official simple-relmgt found. Using local built one."
                    exit 0
                fi
                ''')

                script {
                    officialReleaseFileFound = sh(script: '[ -f ' + releaseCmdPath + ' ]', returnStatus: true)
                    if (officialReleaseFileFound == 0) {
                        releaseStatus = sh(script: releaseCmdPath + ' check', returnStatus: true)
                    } else {
                        releaseCmdPath = "./simple-relmgt"
                    }
                }
            }
        }
        stage('Release PR status') {
            when {
                changeRequest target: 'master'
                expression { return officialReleaseFileFound == 0 }
            }
            steps {
                sh(releaseCmdPath + ' status')
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
                sh('''
                    #!/bin/bash -e
                    source ./build-env.sh
                    go test simple-relmgt simple-relmgt/cmds/draftcmd simple-relmgt/cmds/checkcmd simple-relmgt/cmds/releasecmd simple-relmgt/cmds/statecmd simple-relmgt/cmds/tagcmd'''
                )
            }
        }

        stage('Release PR status from built binary') {
            when {
                changeRequest target: 'master'
                expression { return officialReleaseFileFound != 0 }
            }
            steps {
                script {
                    echo('Using built simple-relmgt...')
                    releaseStatus = sh(script: releaseCmdPath + ' check', returnStatus: true)
                }
            }
        }

        stage('tag it') {
            when {
                branch 'master'
                expression { return releaseStatus == 0 }
            }
            steps {
                sh(releaseCmdPath + ' tag-it') // git tag, push it and create a draft github release
            }
        }

        stage('Deploy') {
            when { 
                branch 'master' 
                expression { return releaseStatus == 0 }
            }
            steps {
                withCredentials([
                    usernamePassword(credentialsId: 'github-jenkins-cred', usernameVariable: 'GITHUB_USER', passwordVariable: 'GITHUB_TOKEN')
                    ]) {
                    sh('''#!/bin/bash -e
                    source ./build-env.sh
                    publish.sh latest''')
                    sh(releaseCmdPath + ' release-it') // release the draft github release
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
