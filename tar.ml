module P = Core_extended.Process
open P.Command_result

type action =
  | Extract
  | Add

let bark args result =
  print_string "FAILED: tar ";
  List.iter (fun x -> Printf.printf "%s " x) args;
  Printf.printf "=> %s %s" result.stdout_tail result.stderr_tail

let perform action ?(on_error = bark) args =
  let flag = match action with
    | Extract -> "x"
    | Add -> "c"
  in
  let prog = "tar" in
  let flags = Printf.sprintf "-Px%sf" flag in
  let args = flags :: args in
  let result = P.run prog args () in
  match result.status with
  | `Exited (x as code) -> begin
      match code with
      | 0 -> ()
      | _ -> on_error args result
    end
  | _ -> on_error args result


let extract = perform Extract
