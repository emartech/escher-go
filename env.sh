
testCaseRepoCloneDirector="`pwd`"
testCasesDirectoryName=".tests"

testCaseDirectoryPath="$testCaseRepoCloneDirector/$testCasesDirectoryName"

if [ ! -d "$testCaseDirectoryPath" ];then
    git clone https://github.com/adamluzsi/escher-test-suite $testCasesDirectoryName
fi

cd $testCasesDirectoryName
git checkout -- .
git pull
cd $testCaseRepoCloneDirector

export TEST_SUITE_PATH="$testCaseDirectoryPath"
