pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                echo "Cloning Repo"
                sh "docker -v"
                checkout scm
            }
        }
        stage('Build') {
            steps {
                echo "Deploying Docker Container"
                sh "docker compose up"
            }
        }
    }

    post {

        always {
            echo "Cleaning workspace"
            cleanWs()
        }

        success {
            echo "Pipeline successful"
        }

        failure {
            echo "Pipeline failed"
        }
    }
}