# Generated by the Protocol Buffers compiler. DO NOT EDIT!
# source: restore.proto
# plugin: grpclib.plugin.main
import abc
import typing

import grpclib.const
import grpclib.client
if typing.TYPE_CHECKING:
    import grpclib.server

import restore_pb2


class RestoreServiceBase(abc.ABC):

    @abc.abstractmethod
    async def Restore(self, stream: 'grpclib.server.Stream[restore_pb2.Answer, restore_pb2.RestoreRequest]') -> None:
        pass

    def __mapping__(self) -> typing.Dict[str, grpclib.const.Handler]:
        return {
            '/restore.RestoreService/Restore': grpclib.const.Handler(
                self.Restore,
                grpclib.const.Cardinality.STREAM_STREAM,
                restore_pb2.Answer,
                restore_pb2.RestoreRequest,
            ),
        }


class RestoreServiceStub:

    def __init__(self, channel: grpclib.client.Channel) -> None:
        self.Restore = grpclib.client.StreamStreamMethod(
            channel,
            '/restore.RestoreService/Restore',
            restore_pb2.Answer,
            restore_pb2.RestoreRequest,
        )
