pipeline {
    agent any

    environment {
        DOCKER_IMAGE = '400034/yourvibes_api_server'
        DOCKER_TAG = 'latest'
        PROD_USER = credentials('PROD_USER')
        PROD_PASSWORD = credentials('PROD_PASSWORD')
        PROD_SERVER = credentials('PROD_SERVER')
        TELEGRAM_BOT_TOKEN = credentials('TELEGRAM_BOT_TOKEN')
        TELEGRAM_CHAT_ID = credentials('TELEGRAM_CHAT_ID')
    }

    stages {
        stage('Clone Repository') {
            steps {
                git branch: 'master', url: 'https://github.com/poin4003/yourVibes_GoApi.git'
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

        stage('Deploy Golang to DEV') {
            steps {
                script {
                    echo 'Clearing server_golang-related images and containers...'
                    sh '''
                        docker container stop yourvibes_api_server || echo "No container named yourvibes_api_server to stop"
                        docker container rm yourvibes_api_server || echo "No container named yourvibes_api_server to remove"
                        docker image rmi ${DOCKER_IMAGE}:${DOCKER_TAG} || echo "No image ${DOCKER_IMAGE}:${DOCKER_TAG} to remove"
                    '''

                    echo 'Deploying to DEV environment...'
                    sh '''
                        docker image pull ${DOCKER_IMAGE}:${DOCKER_TAG}
                        docker network create dev || echo "Network already exists"
                        docker container run -d --rm --name yourvibes_api_server -p 8080:8080 --network dev ${DOCKER_IMAGE}:${DOCKER_TAG}
                    '''
                }
            }
        }

        stage('Deploy to Production on Acer Archlinux server') {
            steps {
                script {
                    echo 'Deploying to Production...'
                    sshScript remote: [
                        host: "${PROD_SERVER}",
                        user: "${PROD_USER}",
                        password: "${PROD_PASSWORD}"
                    ], script: '''
                         docker container stop yourvibes_api_server || echo "No container to stop"
                         docker container rm yourvibes_api_server || echo "No container to remove"
                         docker image rmi ${DOCKER_IMAGE}:${DOCKER_TAG} || echo "No image to remove"
                         docker image pull ${DOCKER_IMAGE}:${DOCKER_TAG}
                         docker container run -d --rm --name yourvibes_api_server -p 8080:8080 ${DOCKER_IMAGE}:${DOCKER_TAG}
                    '''
                }
            }
        }
    }

    post {
        always {
            cleanWs()
        }

        success {
            sendTelegramMessage("✅ Build #${BUILD_NUMBER} was successful! ✅", "${TELEGRAM_BOT_TOKEN}", "${TELEGRAM_CHAT_ID}")
        }

        failure {
            sendTelegramMessage("❌ Build #${BUILD_NUMBER} failed. ❌", "${TELEGRAM_BOT_TOKEN}", "${TELEGRAM_CHAT_ID}")
        }
    }
}

def sendTelegramMessage(String message, String token, String chatId) {
    sh """
    curl -s -X POST https://api.telegram.org/bot${token}/sendMessage \
    -d chat_id=${chatId} \
    -d text="${message}"
    """
}