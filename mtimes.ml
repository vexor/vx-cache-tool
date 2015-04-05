open Big_int
module Time = Core.Time

let fresh_mtime path =
  path, Time.now()

let to_mtime assoc =
  match assoc with
  | `Assoc (path, `Intlit mtime) -> (path, mtime)
  | _ -> assert false

let restore file =
  let json = Yojson.Basic.from_file file in
  let open Yojson.Basic.Util in
  let mtimes = json |> member "mtimes" |> to_assoc |> to_mtime in
  mtim
