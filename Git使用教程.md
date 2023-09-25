# Git

git这个东西，从大二开始接触，一直没有好好花时间搞清楚，现在也不会用，这次下决心一定要搞清楚，顺便写个最适合自己的教程，看别人的教程总是有七零八落与自己的脑子不适配的感觉

## 基本概念

- 工作区： 本地的文件夹，写代码时在工作区写
- 版本库： 将所有版本以树的形式存储下来
- 暂存区： 暂存区是工作区和版本库之间的桥梁，先将工作区交到暂存区，再将暂存区中内容放置入版本库
- 工作区 <==> 暂存区 <==> 版本库，当将暂存区中内容提交时，会在版本库中的HEAD节点后面新增一个节点

## Git开始使用

```bash
git config --global user.name wzp
git config --global user.email xxx@163.com
# 运行上述命令后 .gitconfig文件中可以看到这两个信息
# 其中邮箱最好使用 在云端仓库中账号的邮箱，据说这样在本地中提交后 云端会自动同步（但我试了下并没有）
```

```bash
mkdir project
cd project
git init # 将project文件夹创建成一个git仓库，将出现一个 .git 文件夹
```

```bash
vim readme.txt # 写入一些新内容
git status # 在写入内容后可以看到下面信息
```

![image-20230921212151613](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230921212151613.png)

分析：在master分支下（当前只有一个分支，即master分支），目前没有commits，但是存在未被追踪的文件（Untracked是指当前文件未被加入到暂存区中）

```bash
git add readme.txt # 将文件存入暂存区中，之后才可以将文件commit到版本库中
```

![image-20230921212416543](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230921212416543.png)

分析：status中可以看到暂存区中有新的内容出现，等待被commit到版本库中，并且在暂存区中的文件名以绿名显示

```bash
git rm --cached <file>...  # 可以通过该命令将文件 从暂存区中移除，取消跟踪
```

```bash
git commit #将暂存区中的内容存入到当前分支
git commit -m "add readme.txt" # -m参数可以增加备注 m 即memory，用于记忆该commit内容
```

如果后续修改了代码

![image-20230921212939857](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230921212939857.png)

分析：每次commit完之后暂存区会清空，这时检测到工作区与版本库中内容有所不同，可以通过git add将修改后的文件上传至暂存区

![image-20230922094439214](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922094439214.png)

分析： git add之后再查看status，可以看到这里显示的与第一次将新文件放入暂存区时有所不同

```bash
# 将版本库中不存在的文件放入暂存区，将出现：
git rm --cached <file>...
# 将版本库中已存在的文件的修改后版本放入暂存区后，将出现：
git restore --staged <file>
```

简单理解：restore只是把对当前工作区中文件的修改恢复到修改前的状态（恢复为暂存区中的把呢不能，如果暂存区中为空，就恢复到HEAD所在的commit版本），适用于在写代码时发现当前修改有问题；而rm是指对当前工作区中的某个文件不进行版本管理

```bash
git diff # 查看尚未缓存的改动（我认为可以理解为工作区和版本库中上一次commit的区别，因为此时暂存区中并没有内容）
git diff --cached # 查看已缓存的改动（可以理解为暂存区中与版本库中上一次commit的区别）
git diff HEAD # 查看已缓存的和未缓存的所有改动
git diff --stat # 显示摘要而非整个diff
```

补充一下，对文件的修改也是可以一样要add和commit的

![image-20230922113215405](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922113215405.png)

- 总结：截至这里我们已经学会了基本使用步骤，在工作区中修改代码，使用add命令将内容加入到暂存区，再使用commit命令将内容提交至版本库

## 同一分支中的版本管理

- 我们可以通过如下形式在分支中看到历史版本内容

```bash
# 使用下述两个命令可以查看版本库中当前分支的所有commit记录
git log # 默认包含三类信息：Author Date 备注
git log --pretty=oneline # 只显示 备注信息
```

<img src="C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922104152431.png" alt="image-20230922104152431" style="zoom: 80%;" />

- 如果发现当前版本存在问题，在查看了历史版本信息后可以回滚版本！！！

```bash
git reset --hard HEAD^ # 回滚1个版本
git reset --hard HEAD^^ # 回滚2个版本
# ^ 的数量表示了回滚版本数量
git reset --hard HEAD 版本号
```

![image-20230922105357038](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922105357038.png)

假设有历史版本 1 => 2 => 3 ，当前在版本3，git log可以看到完整的1 => 2 => 3，但是当我们回滚1个版本后，再用git log只能看到 1 => 2，当我们回滚2个版本后，再用git log只能看到 1 。但是！我们仍然使用下面命令看到完整历史版本信息：

```bash
# 回滚并不会删掉已经存在过的版本
git reflog # 显示所有HEAD的移动记录
```

![image-20230922105412768](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922105412768.png)

分析：每一个版本都有一个版本号，是前7位黄名字符，git reset命令也可以利用这个版本号回滚到指定版本，这样避免了中间间隔版本较多时的困扰

```bash
git restore <file> # 将当前工作区内容恢复到暂存区中的内容
git restore --staged <file> # 将当前暂存区中的内容放弃 stage就是暂存区
```

- 总结：这样我们就学会了如何回溯历史版本，在对代码进行一定修改后，保持add并commit的习惯，如果发现后续修改出现错误，还可以回溯到原来的版本！！
- 思考：我们在写代码时，当认为工作区中的修改有价值时，先将其提交至暂存区，当一次代码开发结束，认为暂存区中的修改差不多的时候，可以将暂存区中内容commit为一个版本。

## 云端仓库与本地文件夹协同

- 有多个云端仓库可以选择：Github，GitLab，AC Git

- ssh密钥：创建本地密钥对，将~/.ssh/id_rsa.pub添加到云端仓库中的SSH密钥中

- 如何关联本地仓库与云端仓库？

  - 先有本地Git仓库，再有云端仓库，将本地有内容的Git仓库同步到云端的空仓库中

    - ```bash
      git init # 本地文件夹未构建为Git仓库时需要该操作
      git add .
      git commit -m "first commit"
      git remote add origin git@git.acwing.com:yxc/project.git、
      git push -u origin master # 将本地的master分支推送到origin主机的master分支
      git push <远程主机名> <本地分支名>:<远程分支名>
      git push <远程主机名> <本地分支名> # 本地与远程分支名相同时可省略
      ```

  -  先有云端仓库，从云端仓库抓取内容到本地文件夹(会损失git reflog信息)

    - ```bash
      git clone
      git clone git@git.acwing.com:yxc/project # 本质是scp
      git clone https://git.acwing.com/yxc/project 
      git push -u # 第一次推送commit到远程仓库
      ```

- 在已经建立了云端仓库和本地的联系后，当在本地提交commit后，需要使用push命令将本地commit提交至云端仓库：

  ```bash
  git push -u # 将当前分支合并到远程仓库，第一次带-u参数，后续不需要-u
  ```

  

## 多分支的创建与合并（难点）

- 分支创建与切换：

  新分支的创建位置 可以是任意的，比如说在master分支上已经有多个节点，可以在其任意节点位置创建新分支，新分支会在当前节点处引出一个分叉

  ```bash
  git checkout -b dev # 创建新分支dev并切换到新分支 -b参数指创建分支
  git checkout master # 切换至master分支
  git branch # 查看所有已有分支
  git branch -d dev # 删除分支
  ```

  ![image-20230922133229644](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922133229644.png)

  其中会用*标记当前所处于的分支

- 分支概念与（工作区，暂存区，版本库）之间的关系：

  分支与工作区和暂存区之间是独立开的，add至暂存区中的操作不论在哪个分支下都一样；

  当commit时，会在当前所在分支的HEAD节点后面增加一个节点，即版本库与commit相关

- 多分支：

  多分支的意义在于 修改可以现在非master分支上（如dev）进行，最后将dev分支上的内容合并在master分支上

  ```bash
  git merge dev # 当前在master分支下，将dev分支上的内容合并到当前分支的当前节点上
  ```

  ![image-20230922135647805](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922135647805.png)

  分析：faster forward指快速合并模式，直接将master分支当前节点指向要合并的分支节点（？？？这里没讲清楚）

  合并时候可能会产生冲突： master分支中改了个版本，dev分支中也改了个版本：

  <img src="C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922140526830.png" alt="image-20230922140526830" style="zoom:50%;" />

  那如果将dev合并到master分支就会产生冲突：

  ![image-20230922141154224](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922141154224.png)

  git status查看当前状态：（这里可以想到git status能够查看当前git管理下的任何问题）

  ![image-20230922141245003](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922141245003.png)

  分析：首先第一条信息表示本地仓库比云端仓库领先了2个commits，可以用git push命令将领先的commits同步到云端仓库； 第二条信息表示存在未完成的merge操作，可以处理merge冲突后用commit提交修改，也可以选择 abort终止合并操作； 第三条信息表示两个分支都对readme.txt进行了修改。

  我们可以在master分支下打开readme.txt查看冲突（因为当前操作是把dev合并到master）:

  ![image-20230922142103985](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922142103985.png)

  =====前面是HEAD的内容，后面是dev2的内容，在该文件中可以做任意修改

  ```bash
  git commit -m "fix merge conflicts" # 修改完成后提交需修改
  git log # 查看当前日志信息
  ```

  ![image-20230922142537945](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922142537945.png)

  分析：log按时间顺序记录下了修改，先在dev2上加了”888“，再在master上新增了”999“，最后进行merge处理冲突

  <img src="C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922142852934.png" alt="image-20230922142852934" style="zoom: 50%;" />

  ```bash
  git branch -d dev2 # 因为合并后dev2分支已经没有意义，直接删除即可
  ```

  

- 其他：如果存在多个分支，在其他分支中的多个commit是否进行merge时会发生什么？？？

## 多分支的云端协同（push&pull）

### 本地分支上传到云端

当在本地新建了一个分支并提交了commit之后，需要将该分支内容同步到云端，但如果直接git push会有问题：

<img src="C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922155623982.png" alt="image-20230922155623982" style="zoom: 80%;" />

分析：git push命令将会把当前所在分支的内容传至云端，但是由于云端此时只有master分支，故git找不到对应的传输位置”no upstream branch“

![image-20230922155906367](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922155906367.png)

分析：使用提示信息中的命令将会在云端仓库创建一个新分支dev3并将其与本地分支dev3关联

云端分支内容与本地分支内容具有独立性

```bash
git push -d origin branch_name # 可以删除远程仓库中的分支
```

### 云端分支下载到本地

本地运行下面命令创建新分支

```bash
git checkout -b dev4 
```

```bash
git push --set-upstream origin branch_name # 设置本地的branch_name分支对应远程仓库的branch_name分支
git push --set-upstream-to=origin/branch_name1 branch_name2 # 设置本地的branch_name2分支对应远程仓库的branch_name1分支
```

![image-20230922163203014](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922163203014.png)

绑定之后将远程仓库内容拉下来

```bash
git pull
```

![image-20230922163426357](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922163426357.png)

此时假设本地存在两个分支： ***master***, ***dev4***

假设这样的场景，我们作为项目管理者，其中一个成员将最近的修改提交到了分支***dev4***，我们将***dev4***拉到本地之后，需要将***dev4***中的内容合并到***master***分支中，于是我们执行下述命令

```bash
git checkout master
git merge dev4
# 如果出现冲突，则到产生冲突的具体文件中处理冲突
git branch -d dev4 # 删除本地的dev4分支
git push -d origin dev4 # 删除远程仓库中的dev4分支
git push # 将master分支同步到远程
```

![image-20230922164208801](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20230922164208801.png)

## git stash

当对工作区和暂存区做出了一些修改，但是又不想commit，可以使用该命令

```bash
git stash
git stash apply
git stash drop
git stash pop # 最好在非master分支中弹出
git stash list
```

