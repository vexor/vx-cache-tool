module Sys = Core.Core_sys

type t =
  {
    cacher_dir : string;
    mtime_file : string;
    md5_file   : string;
    fetch_tar  : string;
    push_tar   : string;
  }

let new_config () =
  let cacher_dir =
    match Sys.getenv "CACHER_DIR" with
    | Some path -> path
    | None -> match Sys.getenv "HOME" with
              | Some path -> path ^ "/.cacher"
              | None -> assert false
  in
  {
    cacher_dir = cacher_dir;
    mtime_file = cacher_dir ^ "/mtime.yml";
    md5_file   = cacher_dir ^ "/md5.yml";
    fetch_tar  = cacher_dir ^ "/fetch.tgz";
    push_tar   = cacher_dir ^ "/push.tgz";
  }
