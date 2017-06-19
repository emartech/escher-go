
testCaseRepoCloneDirector="`pwd`"
testCasesDirectoryName=".tests"

testCaseDirectoryPath="$testCaseRepoCloneDirector/$testCasesDirectoryName"

if [ ! -d "$testCaseDirectoryPath" ];then
    git clone "https://github.com/EscherAuth/test-cases" $testCasesDirectoryName
fi

cd $testCasesDirectoryName
git checkout -- .
git pull
cd $testCaseRepoCloneDirector

export TEST_CASES_PATH="$testCaseDirectoryPath"
