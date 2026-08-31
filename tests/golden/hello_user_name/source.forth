variable name_ptr
variable idx
variable ready

: handle_input
  key dup 10 = if
    drop
    0 name_ptr @ idx @ + !
    1 ready !
  else
    name_ptr @ idx @ + !
    idx @ char+ idx !
  then
  iret
;

: puts
  begin
    dup @ 0= not
  while
    dup @ emit
    char+
  repeat
  drop
;

s"                 " name_ptr !
0 idx !
0 ready !
s" What is your name?" puts
10 emit
ei
begin
  ready @
until
di
s" Hello, " puts
name_ptr @ puts
33 emit
