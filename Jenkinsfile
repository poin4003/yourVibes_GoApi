pipeline {
    agent any

    environment {
        DOCKER_IMAGE = '400034/yourvibes_api_server'
        DOCKER_TAG = 'latest'
        PROD_SERVER_PORT = credentials('PROD_SERVER_PORT')
        PROD_SERVER_NAME = credentials('PROD_SERVER_NAME')
        PROD_USER = credentials('PROD_USER')
        PROD_PASSWORD = credentials('PROD_PASSWORD')
        TELEGRAM_BOT_TOKEN = credentials('TELEGRAM_BOT_TOKEN')
        TELEGRAM_CHAT_ID = credentials('TELEGRAM_CHAT_ID')
    }

    stages {
        stage('Clone Repository') {
            steps {
                git branch: 'master', url: 'https://github.com/poin4003/yourVibes_GoApi.git'
            }
        }

        stage('Prepare Config') {
            steps {
                withCredentials([file(credentialsId: 'config_file', variable: 'CONFIG_FILE')]) {
                    sh 'mkdir -p $WORKSPACE/config'
                    sh 'cp $CONFIG_FILE $WORKSPACE/config'
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    echo 'Building Docker image for linux/amd64 platform...'
                    docker.build("${DOCKER_IMAGE}:${DOCKER_TAG}", "--platform linux/amd64 .")
                }
            }
        }

        stage('Run Tests') {
            steps {
                echo 'Running tests...'
            }
        }

        stage('Push to Docker Hub') {
            steps {
                script {
                    docker.withRegistry('https://index.docker.io/v1/', 'docker-hub-credentials') {
                        docker.image("${DOCKER_IMAGE}:${DOCKER_TAG}").push()
                    }
                }
            }
        }

//         stage('Deploy Golang to DEV') {
//             steps {
//                 script {
//                     echo 'Clearing server_golang-related images and containers...'
//                     sh '''
//                         docker container stop yourvibes_api_server || echo "No container named yourvibes_api_server to stop"
//                         docker container rm yourvibes_api_server || echo "No container named yourvibes_api_server to remove"
//                         docker image rm ${DOCKER_IMAGE}:${DOCKER_TAG} || echo "No image ${DOCKER_IMAGE}:${DOCKER_TAG} to remove"
//                     '''
//
//                     echo 'Setting up volume for configuration...'
//                     sh '''
//                         ls -l $WORKSPACE/config
//                         cat $WORKSPACE/config/local.yaml
//                         docker volume create yourvibes_config || echo "Volume yourvibes_config already exists"
//                         docker run --rm -v yourvibes_config:/config --name helper busybox sh -c "mkdir -p /config"
//                         docker cp $WORKSPACE/config/local.yaml helper:/config
//                     '''
//
//                     echo 'Deploying to DEV environment...'
//                     sh '''
//                         docker pull ${DOCKER_IMAGE}:${DOCKER_TAG}
//                         docker run -d --name yourvibes_api_server -p 8080:8080 \
//                         -v yourvibes_config:/config \
//                         ${DOCKER_IMAGE}:${DOCKER_TAG}
//                     '''
//                 }
//             }
//         }

        stage('Deploy to Production on Acer Archlinux server') {
            steps {
                script {
                    echo 'Deploying to Production...'

                    sh '''
                        sshpass -p "${PROD_PASSWORD}" ssh -o StrictHostKeyChecking=no -p "${PROD_SERVER_PORT}" "${PROD_USER}"@${PROD_SERVER_NAME} "
                            echo 'Stopping and removing existing containers and images...'
                            docker container stop yourvibes_api_server || echo 'No container to stop' && \
                            docker container rm yourvibes_api_server || echo 'No container to remove' && \
                            docker image rmi 400034/yourvibes_api_server:latest || echo 'No image to remove'
                        "
                    '''

                    sh '''
                        echo 'Setting up volume for production configuration...'
                        sshpass -p "${PROD_PASSWORD}" ssh -o StrictHostKeyChecking=no -p "${PROD_SERVER_PORT}" "${PROD_USER}"@${PROD_SERVER_NAME} "
                            docker volume create yourvibes_config || echo 'Volume yourvibes_config already exists' && \
                            docker run --rm -v yourvibes_config:/config --name helper busybox sh -c 'mkdir -p /config' && \
                            docker cp ${WORKSPACE}/config/local.yaml helper:/config && \
                        "
                    '''

                    sh '''
                        echo 'Deploying to production server...'
                        sshpass -p "${PROD_PASSWORD}" ssh -o StrictHostKeyChecking=no -p "${PROD_SERVER_PORT}" "${PROD_USER}"@${PROD_SERVER_NAME} "
                            docker pull 400034/yourvibes_api_server:latest && \
                            docker run -d --name yourvibes_api_server -p 8080:8080 \
                            -v yourvibes_config:/config \
                            400034/yourvibes_api_server:latest
                        "
                    '''
                }
            }
        }
    }

    post {
        success {
            cleanWs()
            sendTelegramMessage("✅ Build #${BUILD_NUMBER} was successful! ✅")
        }

        failure {
            cleanWs()
            sendTelegramMessage("❌ Build #${BUILD_NUMBER} failed. ❌")
        }
    }
}

def sendTelegramMessage(String message) {
    withEnv(["MESSAGE=${message}"]) {
        sh '''
        curl -s -X POST https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendMessage \
        -d chat_id=$TELEGRAM_CHAT_ID \
        -d text="$MESSAGE"
        '''
    }
}
