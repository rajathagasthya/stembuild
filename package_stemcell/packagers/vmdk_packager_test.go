package packagers_test

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/package_parameters"
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/packagers"
	"github.com/cloudfoundry-incubator/stembuild/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stemcell", func() {
	var tmpDir string
	var stembuildConfig package_parameters.VmdkPackageParameters
	var c packagers.VmdkPackager

	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		stembuildConfig = package_parameters.VmdkPackageParameters{
			OSVersion: "2012R2",
			Version:   "1200.1",
		}

		c = packagers.VmdkPackager{
			Stop:         make(chan struct{}),
			Debugf:       func(format string, a ...interface{}) {},
			BuildOptions: stembuildConfig,
		}
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	Describe("vmdk", func() {
		Context("valid vmdk file specified", func() {
			It("should be valid", func() {

				vmdk, err := ioutil.TempFile("", "temp.vmdk")
				Expect(err).ToNot(HaveOccurred())
				defer os.Remove(vmdk.Name())

				valid, err := packagers.IsValidVMDK(vmdk.Name())
				Expect(err).To(BeNil())
				Expect(valid).To(BeTrue())
			})
		})
		Context("invalid vmdk file specified", func() {
			It("should be invalid", func() {
				valid, err := packagers.IsValidVMDK(filepath.Join("..", "out", "invalid"))
				Expect(err).To(HaveOccurred())
				Expect(valid).To(BeFalse())
			})
		})
	})

	Describe("CreateImage", func() {

		It("successfully creates an image tarball", func() {
			c.BuildOptions.VMDKFile = filepath.Join("..", "..", "test", "data", "expected.vmdk")
			err := c.CreateImage()
			Expect(err).NotTo(HaveOccurred())

			// the image will be saved to the VmdkPackager's temp directory
			tmpdir, err := c.TempDir()
			Expect(err).NotTo(HaveOccurred())

			outputImagePath := filepath.Join(tmpdir, "image")
			Expect(c.Image).To(Equal(outputImagePath))

			// Make sure the sha1 sum is correct
			h := sha1.New()
			f, err := os.Open(c.Image)
			Expect(err).NotTo(HaveOccurred())

			_, err = io.Copy(h, f)
			Expect(err).NotTo(HaveOccurred())

			actualShasum := fmt.Sprintf("%x", h.Sum(nil))
			Expect(c.Sha1sum).To(Equal(actualShasum))

			// expect the image ova to contain only the following file names
			expectedNames := []string{
				"image.ovf",
				"image.mf",
				"image-disk1.vmdk",
			}

			imageDir, err := helpers.ExtractGzipArchive(c.Image)
			Expect(err).NotTo(HaveOccurred())
			list, err := ioutil.ReadDir(imageDir)
			Expect(err).NotTo(HaveOccurred())

			var names []string
			infos := make(map[string]os.FileInfo)
			for _, fi := range list {
				names = append(names, fi.Name())
				infos[fi.Name()] = fi
			}
			Expect(names).To(ConsistOf(expectedNames))

			// the vmx template should generate an ovf file that
			// does not contain an ethernet section.
			ovf := filepath.Join(imageDir, "image.ovf")
			ovfFile, err := helpers.ReadFile(ovf)
			Expect(err).NotTo(HaveOccurred())
			Expect(ovfFile).NotTo(MatchRegexp(`(?i)ethernet`))
		})
	})

	Describe("CreateManifest", func() {
		It("Creates a manifest correctly", func() {
			expectedManifest := `---
name: bosh-vsphere-esxi-windows1-go_agent
version: 'version'
sha1: sha1sum
operating_system: windows1
cloud_properties:
  infrastructure: vsphere
  hypervisor: esxi
stemcell_formats:
- vsphere-ovf
- vsphere-ova
`
			result := packagers.CreateManifest("1", "version", "sha1sum")
			Expect(result).To(Equal(expectedManifest))
		})
	})
})