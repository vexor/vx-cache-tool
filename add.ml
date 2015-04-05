module Filename = Core.Core_filename
module Unix = Core.Core_unix

let add_path fetch_tar path =
  let path = Filename.realpath path in
  Printf.printf "adding %s to cache" path;
  Unix.mkdir_p path;
  Tar.extract [fetch_tar; path] ~on_error:(
                fun _ _ -> Printf.printf "%s isn't yet cached" path
              );
  Mtimes.fresh_mtime path

let run dirs () =
  let cfg = Config.new_config () in
  let mtimes = Mtimes.restore cfg.mtimes_file in
  List.map (add_path cfg.fetch_tar) dirs
  |> Mtimes.refresh mtimes
  |> Mtimes.store
