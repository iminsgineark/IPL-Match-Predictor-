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
        stage('Image Build') {
            steps {
                sh "docker build -t ${IMAGE_NAME} -f ${DOCKERFILE_PATH}/Dockerfile ${DOCKERFILE_PATH}"
            }
        }
        // stage('pytest tests'){
        //     steps{
        //         sh "pytests tests/"
        //     }
        // }        
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
