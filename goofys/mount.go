package goofys

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseutil"
)

func Mount(bucket, mountPoint string) (mfs *fuse.MountedFileSystem, err error) {
	awsConfig := &aws.Config{
		Region: aws.String("us-east-1"),
	}
	awsConfig.S3ForcePathStyle = aws.Bool(true)
	dur, _ := time.ParseDuration("1m")
	flags := &FlagStorage{
		MountOptions: make(map[string]string),
		DirMode:      os.FileMode(600),
		FileMode:     os.FileMode(600),
		Uid:          uint32(99),
		Gid:          uint32(99),
		StorageClass: "STANDARD",
		StatCacheTTL: dur,
		TypeCacheTTL: dur,
		DebugFuse:    false,
		DebugS3:      false,
		Foreground:   false,
	}

	fs := NewGoofys(bucket, awsConfig, flags)
	server := fuseutil.NewFileSystemServer(fs)
	mountCfg := &fuse.MountConfig{
		FSName:                  bucket,
		Options:                 flags.MountOptions,
		DisableWritebackCaching: true,
	}
	mfs, err = fuse.Mount(mountPoint, server, mountCfg)
	return
}
