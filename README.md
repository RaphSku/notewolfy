# notewolfy
## Description
notewolfy simplifies note-taking and organization into a straightforward tree-like structure. Currently, notewolfy supports note-taking through Markdown files via Vim. If you're unfamiliar with Vim, I recommend learning it before using notewolfy. Take notes of something noteworthy.

## How to get it
You can get notewolfy with the following command:
```bash
go install github.com/RaphSku/notewolfy@latest
```

## How to use it
notewolfy is a console application, so open your terminal and type in:
```bash
notewolfy
```
You will see a prompt where you can start typing:
```bash
>>> create workspace example ~/example
```
This will create a new workspace for you named `example` at the path `~/example`. Before we start taking notes, consider organizing your workspace further. For instance, if you're researching something, you might want to create a `research` node within your workspace.
```bash
>>> create node research
```
Now, move to your new node with
```bash
>>> goto research
```
And let us create a Markdown file.
```bash
>>> create md research_topic_a
```
Edit your new Markdown file and write something into it.
```bash
>>> edit research_topic_a
```
You will see that Vim opens and will let you edit your Markdown file. Furthermore, you have seen how to navigate forward but how do we get back? Well, just go back.
```bash
>>> goback
```
By the way, at any time you can use
```bash
>>> ls
```
to see information about the current node you are on and which Markdown files reside there and you can use
```bash
>>> ls ws
```
to see information about all your workspaces, namely the name and the path where they reside. This is useful if you forget the path to your workspace.
If one workspace does not suffice, just create another workspace
```bash
>>> create workspace example2 ~/example2
```
You will not need to switch to the new workspace since this is done automatically. But if you want to open a particular workspace, you can simply use
```bash
>>> open example2
```
Once you are finished, you can close notewolfy by pressing either keys: "Esc", "Ctrl+C" or type in the following commands: "quit" or "exit".

If you want to delete Markdown files, you can do this with
```bash
>>> delete md <markdownFileName>
```
, leave the extension .md out when you specify the name. You can only delete Markdown files on the node that you are currently on. You can also delete a node with
```bash
>>> delete node <nodeName>
```
but you need to go to the parent node to delete the child node. If you have deleted all Markdown files and nodes, you can also delete the workspace.
```bash
>>> delete workspace <workspaceName>
```
A bulk delete is currently not supported but if you want to delete the whole workspace without going over every node and Markdown file, you can simply delete it via the file explorer or terminal. You also need to remove the workspace metadata in the `.notewolfy` metadata file that was created in your home directory. It is JSON encoded, so just remove the workspace entry under workspaces.

If you need help with a command, try to use
```bash
help create workspace
```

If you want to see the version of notewolfy that you are using, just use the following command
```bash
notewolfy version
```
or inside of notewolfy
```bash
>>> version
```
there is also the version command available.

## Supported OS
Only UNIX operating systems are supported. Sorry, no Windows for the time being.
