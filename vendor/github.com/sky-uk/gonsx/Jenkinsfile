#!groovy

project_name = 'gonsx'
project_owner = 'sky-uk'
project_github_token = '7eec88ceb08431126182c7ce85631cb45bd0f98a'

project_src_path = "github.com/${project_owner}/${project_name}"
git_credentials_id = '9be96924-ccbc-4f9e-a07c-18818fff868c'
git_helper_credentials_id = 'paas-jenkins-pipelines-deploy-key'

version_file = 'VERSION'
major_version = null
minor_version = null
patch_version = null

docker_image = "paas/golang-img:0.8.8"

// helpers
gitHelper = null
shellHelper = null
goHelper = null
slackHelper = null

slackChannel = '#ott-paas'

loadHelpers()

slackHelper.notificationWrapper(slackChannel, currentBuild, env, true) {
    node {
        wrap([$class: 'TimestamperBuildWrapper']) {
            wrap([$class: 'AnsiColorBuildWrapper']) {
                stage 'checkout'
                deleteDir()
                git_branch = env.BRANCH_NAME
                checkout scm
                gitHelper.prepareGit('paas-jenkins', 'paas-jenkins@jenkins.paas.int.ovp.bskyb.com')

                stage 'version'
                if (autoincVersion()) {
                    writeFile file: version_file, text: version()
                    gitHelper.commit(version_file, "bumping to: ${version()}")
                }

                echo "Starting pipeline for project: [${project_name}], branch: [${git_branch}], version: [${version()}]"

                stage 'lint'
                inContainer {
                    goHelper.goLint(project_src_path)
                }

                stage 'format'
                inContainer {
                    goHelper.goFmt(project_src_path)
                }

                stage 'vet'
                inContainer {
                    goHelper.goVet(project_src_path)
                }

                stage 'build'
                inContainer {
                    goHelper.goBuild(project_src_path)
                }

                stage 'test'
                inContainer {
                    goHelper.goTest(project_src_path)
                }

                stage 'coverage'
                inContainer {
                    goHelper.goCoverage(project_src_path)
                }
            }
        }
    }
    // we only release from master
    if (git_branch == 'master' && !gitHelper.isLastCommitFromUser('paas-jenkins')) {
        stage 'release'
        def approved = true
        timeout(time: 2, unit: 'HOURS') {
            try {
                input message: "Release ${project_name} ${version()} ?"
            } catch (InterruptedException _x) {
                echo "Releasing not approved in time!"
                approved = false
            }
        }

        if (approved) {
            echo "Release has been approved!"
            node() {
                gitHelper.tag(version(), "Jenkins ${env.JOB_NAME} ${env.BUILD_DISPLAY_NAME}")
                gitHelper.push(git_credentials_id, git_branch)

                echo "Creating GitHub Release v${version()}"
                github_release_response = gitHelper.createGitHubRelease(project_github_token, project_owner, project_name, version(), git_branch)
//              FIXME: this is not working yet
//              echo "Attaching artifacts to GitHub Release v${version()}"
//              gitHelper.uploadToGitHubRelease(project_github_token, project_owner, project_name, github_release_response.id, "${pwd()}/coverage.html", 'application/html')
            }
            currentBuild.description = "Released ${version()}"
        }
    }
}

def loadHelpers() {
    fileLoader.withGit('git@github.com:sky-uk/paas-jenkins-pipelines.git', 'master', git_helper_credentials_id, '') {
        this.gitHelper = fileLoader.load('lib/helpers/git')
        this.shellHelper = fileLoader.load('lib/helpers/shell')
        this.goHelper = fileLoader.load('lib/helpers/go')
        this.slackHelper = fileLoader.load('lib/helpers/slack')
    }
}

def autoincVersion() {
    current_version = readFile("${pwd()}/${this.version_file}").trim().tokenize(".")
    setVersion(current_version[0], current_version[1], current_version[2])

    if(gitHelper.checkIfTagExists(version())) {
        this.patch_version++
        if(gitHelper.checkIfTagExists(version())) {
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
    this.major_version = major.toInteger()
    this.minor_version = minor.toInteger()
    this.patch_version = patch.toInteger()
}

def inContainer(Closure body) {
    docker.image(this.docker_image).inside("-v ${pwd()}:/paas/go/src/${project_src_path} -v ${System.getProperty('java.io.tmpdir')}:${System.getProperty('java.io.tmpdir')}") {
        body()
    }
}
