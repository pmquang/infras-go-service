pipeline {
  agent {
    docker {
      image "${env.AWS_ACCOUNT_ID}.dkr.ecr.ap-southeast-1.amazonaws.com/internal/terraform:0.12.8-alpine-3.10-awscli-1.16.279"
      registryUrl "https://${env.AWS_ACCOUNT_ID}.dkr.ecr.ap-southeast-1.amazonaws.com"
      args '--privileged --volume $DOCKER_CONFIG/:/root/.docker/'
    }
  }

  environment {
    DOCKER_REGISTRY = "${env.AWS_ACCOUNT_ID}.dkr.ecr.ap-southeast-1.amazonaws.com"
    SERVICE_NAME = "infras-go-service"
    COMPONENT = "backend"
    VERSION = sh(script: """cat ./VERSION | tail -1""", returnStdout: true).trim()
    BUILD_IMAGE = "${env.COMPONENT}/${env.SERVICE_NAME}:build-base"
    SERVICE_IMAGE = "${env.COMPONENT}/${env.SERVICE_NAME}:${env.VERSION}"

  }

  stages {
    stage('Prepare') {
      steps {
        sh(label: 'Pull latest build-base Docker image to be served as cache sources of later Docker builds',
           script: 'docker pull $DOCKER_REGISTRY/$BUILD_IMAGE || true')
        sh(label: 'Build build-base Docker image to be served as base environment of later stages',
           script: '''
             set +x
             docker build -t $BUILD_IMAGE \
                    --target build-base \
                    --cache-from $DOCKER_REGISTRY/$BUILD_IMAGE .
           ''')
      }
    }
    stage('Cache and Build') {
      failFast true
      parallel {
        stage('Cache') {
          steps {
            sh(label: 'Tag build-base Docker image as preparation for pushing to ECR',
               script: 'docker tag $BUILD_IMAGE $DOCKER_REGISTRY/$BUILD_IMAGE')
            sh(label: 'Push build-base Docker image to ECR for usages as Docker builds\' cache sources in later pipelines',
               script: 'docker push $DOCKER_REGISTRY/$BUILD_IMAGE')
          }
        }
        stage('Build') {
          steps {
            sh(label: 'Build service Docker image',
               script: '''
                 set +x
                 docker build -t $SERVICE_IMAGE \
                       --cache-from $BUILD_IMAGE .
               ''')
          }
        }
      }
    }

    stage('Deploy Docker Image') {
      steps {
        sh(label: 'Tag service Docker image as preparation for pushing to ECR',
           script: 'docker tag $SERVICE_IMAGE $DOCKER_REGISTRY/$SERVICE_IMAGE')
        sh(label: 'Push service Docker image to ECR',
           script: 'docker push $DOCKER_REGISTRY/$SERVICE_IMAGE')
      }
    }
  }
}
