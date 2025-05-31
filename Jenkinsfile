pipeline {
	agent any

	stages {
		stage('Checkout'){
			steps {
				checkout scm
			}
		}	
 		stage('Build and Run') {
			steps {
				sh '''
					go version
					go build -o api
					./api
				  '''
				}
		}
	}
}
