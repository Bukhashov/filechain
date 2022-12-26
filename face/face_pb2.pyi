from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ComparisonRequest(_message.Message):
    __slots__ = ["forCheckImage", "forCheckMetadata", "originalImage", "originalMetadata"]
    FORCHECKIMAGE_FIELD_NUMBER: _ClassVar[int]
    FORCHECKMETADATA_FIELD_NUMBER: _ClassVar[int]
    ORIGINALIMAGE_FIELD_NUMBER: _ClassVar[int]
    ORIGINALMETADATA_FIELD_NUMBER: _ClassVar[int]
    forCheckImage: bytes
    forCheckMetadata: Metadata
    originalImage: bytes
    originalMetadata: Metadata
    def __init__(self, originalMetadata: _Optional[_Union[Metadata, _Mapping]] = ..., originalImage: _Optional[bytes] = ..., forCheckMetadata: _Optional[_Union[Metadata, _Mapping]] = ..., forCheckImage: _Optional[bytes] = ...) -> None: ...

class ComparisonRespons(_message.Message):
    __slots__ = ["coincidences"]
    COINCIDENCES_FIELD_NUMBER: _ClassVar[int]
    coincidences: bool
    def __init__(self, coincidences: bool = ...) -> None: ...

class FindRequest(_message.Message):
    __slots__ = ["image", "metadata"]
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    image: bytes
    metadata: Metadata
    def __init__(self, metadata: _Optional[_Union[Metadata, _Mapping]] = ..., image: _Optional[bytes] = ...) -> None: ...

class FindRespons(_message.Message):
    __slots__ = ["total"]
    TOTAL_FIELD_NUMBER: _ClassVar[int]
    total: int
    def __init__(self, total: _Optional[int] = ...) -> None: ...

class Metadata(_message.Message):
    __slots__ = ["extension", "filename"]
    EXTENSION_FIELD_NUMBER: _ClassVar[int]
    FILENAME_FIELD_NUMBER: _ClassVar[int]
    extension: str
    filename: str
    def __init__(self, filename: _Optional[str] = ..., extension: _Optional[str] = ...) -> None: ...
