pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                echo "Cloning Repo"
                checkout scm
            }
        }
        stage('Build') {
            steps {
                echo "Deploying Docker Container"
                sh "docker compose up -d"
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