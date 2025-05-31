pipeline {
	agent any
	tools{
		Docker: 'docker'
	}
	environment{
		IMAGE_NAME = "go-api"
		CONTAINER_NAME = "api-container"
		IMAGE_TAG = "${env.BUILD_NUMBER}"
		}
	

	stages {
		stage('Checkout'){
			steps {
				checkout scm
			}
		}
		stage('Build Docker Image') {
			steps {
				sh """
				echo " Building Docker Image: ${IMAGE_NAME}:${IMAGE_TAG}"
				docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .
				"""
			}
		}
		stage('Stop & Remove Old Container') {
			steps {
				sh """
					if [ \$(docker ps -q -f name=${CONTAINER_NAME}) ]; then
						echo "Stopping container ${CONTAINER_NAME}"
						docker stop ${CONTAINER_NAME}
						docker rm ${CONTAINER_NAME}
					else
						echo "No pre-existing container"
					fi
				  """
				}
		}
		stage('Run new container') {
			steps {
				sh """
					echo "Running new container ${CONTAINER_NAME} from image ${IMAGE_NAME}"
					docker run -d \\
					 --name ${CONTAINER_NAME} \\
					-p 8081:8081 \\
					${IMAGE_NAME}:${IMAGE_TAG}
				  """
			}
		}
	}
}
