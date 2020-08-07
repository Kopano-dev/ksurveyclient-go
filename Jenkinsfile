#!/usr/bin/env groovy

pipeline {
	agent {
		docker {
			image 'golang:1.13'
			args '-u 0'
		}
	}
	environment {
		GOBIN = '/usr/local/bin'
	}
	stages {
		stage('Bootstrap') {
			steps {
				echo 'Bootstrapping..'
				sh 'cd / && go get -v golang.org/x/lint/golint'
				sh 'cd / && go get -v github.com/tebeka/go2xunit'
				sh 'cd / && go get -v github.com/axw/gocov/gocov'
				sh 'cd / && go get -v github.com/AlekSi/gocov-xml'
				sh 'go version'
			}
		}
		stage('Lint') {
			steps {
				echo 'Linting..'
				sh 'golint ./... | tee golint.txt || true'
				sh 'go vet ./... | tee govet.txt || true'
			}
		}
		stage('Test') {
			steps {
				echo 'Testing..'
				sh 'go test -v -covermode=atomic -coverprofile=coverage.out | tee tests.output'
				sh 'go2xunit -fail -input tests.output -output tests.xml'
			}
		}
		stage('Coverage') {
			steps {
				echo 'Coverage..'
				sh 'mkdir -p ./test/reports'
				sh 'go tool cover -html=coverage.out -o test/reports/coverage.html'
				sh 'gocov convert coverage.out | gocov-xml > coverage.xml'
			}
		}
	}
	post {
		always {
			junit allowEmptyResults: true, testResults: 'tests.xml'
			recordIssues qualityGates: [[threshold: 100, type: 'TOTAL', unstable: true]], tools: [goVet(pattern: 'govet.txt'), goLint(pattern: 'golint.txt')]
			publishHTML([allowMissing: true, alwaysLinkToLastBuild: true, keepAll: true, reportDir: 'test/reports', reportFiles: 'coverage.html', reportName: 'Go Coverage Report HTML', reportTitles: ''])
			step([$class: 'CoberturaPublisher', autoUpdateHealth: false, autoUpdateStability: false, coberturaReportFile: 'coverage.xml', failUnhealthy: false, failUnstable: false, maxNumberOfBuilds: 0, onlyStable: false, sourceEncoding: 'ASCII', zoomCoverageChart: false])
			cleanWs()
		}
	}
}
