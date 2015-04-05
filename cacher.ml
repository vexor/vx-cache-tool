open Core.Std

let add =
  Command.basic
    ~summary:"Add directories to the cache"
    Command.Spec.(
    empty
      +> anon (sequence ("directory name" %: file))
    )
    Add.run

let command = Command.group ~summary: "Manipulate Vexor cache"
                            ["add", add]

let () = Command.run
           ~version:"v0.0.1"
           ~build_info:"VX-OCaml-experimental"
           command
