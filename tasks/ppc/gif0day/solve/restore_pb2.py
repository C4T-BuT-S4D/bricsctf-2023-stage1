# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: restore.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\rrestore.proto\x12\x07restore\"\x18\n\x06\x41nswer\x12\x0e\n\x06\x61nswer\x18\x01 \x01(\t\"9\n\x0eRestoreRequest\x12\x0c\n\x04hash\x18\x01 \x01(\t\x12\x0c\n\x04salt\x18\x02 \x01(\t\x12\x0b\n\x03gif\x18\x03 \x01(\x0c\x32K\n\x0eRestoreService\x12\x39\n\x07Restore\x12\x0f.restore.Answer\x1a\x17.restore.RestoreRequest\"\x00(\x01\x30\x01\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'restore_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  _globals['_ANSWER']._serialized_start=26
  _globals['_ANSWER']._serialized_end=50
  _globals['_RESTOREREQUEST']._serialized_start=52
  _globals['_RESTOREREQUEST']._serialized_end=109
  _globals['_RESTORESERVICE']._serialized_start=111
  _globals['_RESTORESERVICE']._serialized_end=186
# @@protoc_insertion_point(module_scope)
