pipeline {
    agent any

    stages {

        stage('Test Gudang') {
            steps {
                dir('TUBES/logistik-microservice-final/gudang-service') {
                    bat 'go test ./...'
                }
            }
        }

        stage('Test Courier') {
            steps {
                dir('TUBES/logistik-microservice-final/courier-service') {
                    bat 'go test ./...'
                }
            }
        }

        stage('Build Images') {
            steps {
                dir('TUBES/logistik-microservice-final') {
                    bat 'docker compose build'
                }
            }
        }

        stage('Run Containers') {
            steps {
                dir('TUBES/logistik-microservice-final') {
                    bat 'docker compose up -d'
                }
            }
        }

        stage('Verify') {
            steps {
                bat 'curl http://localhost:8085/packages'
                bat 'curl http://localhost:8086/deliveries'
            }
        }
    }
}