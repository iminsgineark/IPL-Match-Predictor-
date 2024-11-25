pipeline {
    agent any
    environment {
        IMAGE_NAME = "ipl_predictor_app"
        DOCKERFILE_PATH = "."
    }
    stages {
        stage('checkout') {
            steps {
                git branch: 'main', url: 'https://github.com/iminsgineark/IPL-Match-Predictor-'
            }
        }

         stage('Install dependencies') {
            steps {
                script {
                    sh '''
                    if ! dpkg -l | grep -q python3.11-venv; then
                        apt-get update
                        apt-get install -y python3.11-venv
                    fi
                    '''
                }
            }
        }
        stage('Create Virtual Environment') {
            steps {
                script {
                    sh 'python3 -m venv venv'
                }
            }
        }
        stage('Install Python Dependencies') {
            steps {
                script {
                     sh '''
                     . venv/bin/activate
                     pip install -r requirements.txt
                     '''
                }
            }
        }

        stage("bandit scan") {
            steps{
                script {
                    sh '''
                    . venv/bin/activate
                    pip install bandit
                    bandit -r . -f json -o bandit_report.json || true
                    '''
                }
            }
        }

        stage('Image Build') {
            steps {
                sh "docker build -t ${IMAGE_NAME} -f ${DOCKERFILE_PATH}/Dockerfile ${DOCKERFILE_PATH}"
            }
        }
        stage('Trivy Scan') {
            steps {
                sh "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy:latest image ${IMAGE_NAME}:latest --format json --output /tmp/trivy_report.json"
            }
        }
        stage('Container Run') {
            steps {
                sh "docker run -d -p 8501:8501 ${IMAGE_NAME}"
            }
        }
    }
}
