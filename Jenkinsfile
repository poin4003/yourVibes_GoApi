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

        stage('Deploy to Production on Acer Archlinux server') {
            steps {
                script {
                    echo 'Deploying to Production...'

                    sh '''
                        echo 'Stopping and removing existing container...'
                        sshpass -p "${PROD_PASSWORD}" ssh -o StrictHostKeyChecking=no -p "${PROD_SERVER_PORT}" "${PROD_USER}"@${PROD_SERVER_NAME} "
                            docker stop yourvibes_api_server || echo 'Container not running' && \
                            docker rm yourvibes_api_server || echo 'Container not found'
                        "
                    '''

                    sh '''
                        echo 'Removing old Docker image...'
                        sshpass -p "${PROD_PASSWORD}" ssh -o StrictHostKeyChecking=no -p "${PROD_SERVER_PORT}" "${PROD_USER}"@${PROD_SERVER_NAME} "
                            docker rmi 400034/yourvibes_api_server:latest || echo 'Image not found'
                        "
                    '''

                    sh '''
                        echo 'Copying prod.yaml to production server...'
                        sshpass -p "${PROD_PASSWORD}" scp -P "${PROD_SERVER_PORT}" \
                        ${WORKSPACE}/config/prod.yaml \
                        "${PROD_USER}"@${PROD_SERVER_NAME}:/home/pchuy/documents/yourVibes_GoApi/config/
                    '''

                    sh '''
                        echo 'Setting up Docker volume for production configuration...'
                        sshpass -p "${PROD_PASSWORD}" ssh -o StrictHostKeyChecking=no -p "${PROD_SERVER_PORT}" "${PROD_USER}"@${PROD_SERVER_NAME} "
                            docker volume create yourvibes_config || echo 'Volume yourvibes_config already exists' && \
                            docker run --rm -v yourvibes_config:/config -v /home/pchuy/documents/yourVibes_GoApi/config:/host busybox sh -c 'cp /host/prod.yaml /config/prod.yaml'
                        "
                    '''

                    sh '''
                        echo 'Deploying application to production server...'
                        sshpass -p "${PROD_PASSWORD}" ssh -o StrictHostKeyChecking=no -p "${PROD_SERVER_PORT}" "${PROD_USER}"@${PROD_SERVER_NAME} "
                            docker pull 400034/yourvibes_api_server:latest && \
                            docker network connect yourvibes_goapi_default yourvibes_api_server || echo 'Network already connected' && \
                            docker run -d --name yourvibes_api_server -p 8080:8080 \
                                -e YOURVIBES_SERVER_CONFIG_FILE=prod \
                                -v yourvibes_config:/config \
                                -v /etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt:ro \
                                -v yourvibes_goapi_yourvibes_storage:/storages \
                                -v yourvibes_goapi_tmp_volume:/tmp \
                                --dns=8.8.8.8 --dns=8.8.4.4 \
                                --network yourvibes_goapi_default \
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
