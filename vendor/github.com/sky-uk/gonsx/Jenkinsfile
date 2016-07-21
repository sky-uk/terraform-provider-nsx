#!groovy
import java.security.MessageDigest;

project_name = 'GoNSX'
project_github_url = 'git@github.com:sky-uk/gonsx.git'
project_src_path = 'github.com/sky-uk/gonsx'
git_credentials_id = '9be96924-ccbc-4f9e-a07c-18818fff868c'

version_file = 'VERSION'
major_version = null
minor_version = null
patch_version = null

docker_image = "paas/golang-img:0.0.1"

node {
    wrap([$class: 'TimestamperBuildWrapper']) {
        wrap([$class: 'AnsiColorBuildWrapper']) {
            ws {
                stage 'checkout'
                deleteDir()
                git_branch = env.BRANCH_NAME
                echo "Checking out branch [${git_branch}]"
                git branch: git_branch, url: project_github_url, credentialsId: git_credentials_id

                stage 'version'
                if(autoincVersion()) {
                    writeFile file: version_file, text: version()
                    commit(version_file, "bumping to: ${version()}")
                }

                echo "Starting pipeline for project: [${project_name}], branch: [${git_branch}], version: [${version()}]"

                stage 'lint'
                echo "Running Go lint"
                inContainer {
                    sh "go get github.com/golang/lint/golint"
                    sh "\$GOPATH/bin/golint -set_exit_status ${project_src_path}/..."
                }

                stage 'format'
                echo "Verifying source code format with gofmt"
                inContainer {
                    shOut = enhancedSh("gofmt -d -e \$GOPATH/src/${project_src_path}")
                    if(shOut[0] != 0 || shOut[1] != "") {
                        echo "Exit code: ${shOut[0]}"
                        echo "Stdout: ${shOut[1]}"
                        echo "Stderr: ${shOut[2]}"
                        error "gofmt failed!"
                    }
                }

                stage 'vet'
                echo "Running Go vet to find potential issues"
                inContainer {
                    sh "go vet ${project_src_path}/..."
                }

                stage 'build'
                echo "Building ${project_name} ${version()}"
                inContainer {
                    sh "go build ${project_src_path}/..."
                }

                stage 'test'
                echo "Running unit tests"
                inContainer {
                    sh "go get github.com/stretchr/testify/assert"
                    sh "go test ${project_src_path}/..."
                }

                stage 'coverage'
                echo "Generating coverage report"
                inContainer {
                    sh "go get github.com/stretchr/testify/assert"
                    // JaCoCo cobertura output (cobertura plugin not supported in jenkins pipelines yet!
                    sh "go get github.com/axw/gocov"
                    sh "go get github.com/AlekSi/gocov-xml"
                    sh "go get github.com/matm/gocov-html"

                    sh "\$GOPATH/bin/gocov test ${project_src_path}/... | \$GOPATH/bin/gocov-xml > ${pwd()}/coverage-gocov.xml"
                    sh "\$GOPATH/bin/gocov test ${project_src_path}/... | \$GOPATH/bin/gocov-html > ${pwd()}/coverage-gocov.html"
                    publishHTML (target: [
                            allowMissing: false,
                            alwaysLinkToLastBuild: false,
                            keepAll: true,
                            reportDir: pwd(),
                            reportFiles: 'coverage-gocov.html',
                            reportName: "Cobertura Report"
                    ])

                    // TODO: go test --coverprofile doesn't support multiple pacakges yet, so we use this to merge them
                    sh "go get github.com/go-playground/overalls"
                    sh "cd \$GOPATH/src/${project_src_path} && \$GOPATH/bin/overalls -project=${project_src_path} -covermode=count -debug"
                    sh "go tool cover -html=\$GOPATH/src/${project_src_path}/overalls.coverprofile -o coverage.html"

                    publishHTML (target: [
                            allowMissing: false,
                            alwaysLinkToLastBuild: false,
                            keepAll: true,
                            reportDir: '.',
                            reportFiles: 'coverage.html',
                            reportName: "Coverage Report"
                    ])
                }

                // we only release from master
                if(git_branch == 'master') {
                    stage 'release'
                    tag(version(), "Jenkins ${env.JOB_NAME} ${env.BUILD_DISPLAY_NAME}")
                    push(git_credentials_id, git_branch)
                }
            }
        }
    }

}

def autoincVersion() {
    current_version = readFile("${pwd()}/${this.version_file}").trim().tokenize(".")
    setVersion(current_version[0], current_version[1], current_version[2])

    if(checkIfTagExists(version())) {
        this.patch_version++
        if(checkIfTagExists(version())) {
            error "Next patch version (${version()}) already exists!"
        }
        return true
    }
    return false
}

def version() {
    return "${this.major_version}.${this.minor_version}.${this.patch_version}"
}

def setVersion(major, minor, patch) {
    this.major_version = major
    this.minor_version = minor
    this.patch_version = patch
}

// git workarounds as there is no gitpublisher plugin yet
def commit(String filesToAdd, String comment) {
    sh("git add ${filesToAdd}")
    sh("git commit -m '${comment}'")
}

def tag(String tag, String comment) {
    sh("git tag -a -f -m '${comment}' ${tag}")

}

def push(String credentials, String branch) {
    sshagent([credentials]) {
        sh("git push origin ${branch}")
        sh("git push origin ${branch} --tags")
    }
}

def checkIfTagExists(tag) {
    echo "Checking if tag ${tag} exists"
    shOut = enhancedSh("git rev-parse -q --verify \"refs/tags/${tag}\"")
    if(shOut[0] == 0) return true
    return false
}

def inContainer(Closure body) {
    docker.image(this.docker_image).inside("-v ${pwd()}:/paas/go/src/${project_src_path} -v ${System.getProperty('java.io.tmpdir')}:${System.getProperty('java.io.tmpdir')}") {
        body()
    }
}

// FIXME: this function is very hacky... but the "sh" step is very limited atm
def enhancedSh(command) {
    def generateMD5 = generateMD5(command)
    def tmpDir = "${System.getProperty('java.io.tmpdir')}/jenkins-enhancedsh"
    new File(tmpDir).mkdirs()
    def filesPrefix = "${tmpDir}/sh-${generateMD5}"
    def commandFilePath = ("${filesPrefix}-command.txt")
    writeFile file: commandFilePath, text: command

    def exitCodeFilePath = ("${filesPrefix}-exitCode.txt")
    def stdoutFilePath = ("${filesPrefix}-stdout.txt")
    def stderrFilePath = ("${filesPrefix}-stderr.txt")

    def exitCodeFile = new File(exitCodeFilePath)
    def stdoutFile = new File(stdoutFilePath)
    def stderrFile = new File(stderrFilePath)

    exitCodeFile.deleteOnExit()
    stdoutFile.deleteOnExit()
    stderrFile.deleteOnExit()

    echo "Executing [${command}], output to ${filesPrefix}-*.txt"
    sh "set +e; ${command} 1>${stdoutFile.getAbsolutePath()} 2>${stderrFile.getAbsolutePath()} ; echo \$? > ${exitCodeFile.getAbsolutePath()} "
    int exitCode = readFile(exitCodeFile.getAbsolutePath()).trim().toInteger()
    def stdout = readFile(stdoutFile.getAbsolutePath()).trim()
    def stderr = readFile(stderrFile.getAbsolutePath()).trim()

    return [exitCode, stdout, stderr]
}

def generateMD5(String s){
    MessageDigest.getInstance("MD5").digest(s.bytes).encodeHex().toString()
}