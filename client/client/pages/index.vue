<template>
  <div>
    <v-progress-linear
      v-if="fileUploadStart"
      v-model="progress"
      color="light-blue"
      height="25"
    >
      <strong>{{ progress }} %</strong>
    </v-progress-linear>

    <v-row no-gutters justify="center" align="center">
      <v-col cols="8">
        <v-file-input
          show-size
          label="File input"
          v-model="file"
          @change="fileChange"
        ></v-file-input>
      </v-col>
      <br />

      <v-col cols="4" class="pl-2">
        <v-btn color="success" dark small @click="startMultiPartUpload">
          Upload
          <v-icon right dark>mdi-cloud-upload</v-icon>
        </v-btn>
      </v-col>
    </v-row>

    start byte : {{ currentChunkStartByte }} <br />
    end byte : {{ currentChunkFinalByte }} <br />
    size byte : {{ currentChunkFinalByte - currentChunkStartByte }} <br />
    size mb : {{ (currentChunkFinalByte - currentChunkStartByte) / 1000000 }}

    <v-alert v-if="message" border="left" color="blue-grey" dark>
      {{ message }}
    </v-alert>
  </div>
</template>

<script>
export default {
  layout: "",
  data() {
    return {
      message: "",
      file: "",

      currentChunkStartByte: 0,
      currentChunkFinalByte: 6 * 1000000, // default if file size greater then 5 mb
      chunkSize: 6 * 1000000, //  chunk upload size
      minimumChunk: 6 * 1000000, // minimum chunk uppload size
      fileUploadStart: false,
      partNumber: 1,
    };
  },
  computed: {
    progress() {
      return (this.currentChunkFinalByte / this.file.size) * 100;
    },
    uploadChunkSize() {
      return this.currentChunkFinalByte - this.currentChunkFinalByte;
    },
  },

  methods: {
    fileChange() {
      this.currentChunkFinalByte =
        this.chunkSize > this.file.size ? this.file.size : this.chunkSize;
    },
    async startMultiPartUpload() {
      this.fileUploadStart = true;
      const resp = await this.$axios.$post(
        "https://upload.feblic.com/upload/initilize",
        {
          fileType: this.file.name.split(".").pop(),
        }
      );

      this.$axios.setHeader("X-Upload-Id", resp.uploadId);
      this.uploadChunk();
    },
    async uploadChunk() {
      console.log(this.file);

      const bodyFormData = new FormData();

      this.$axios.setHeader("Content-Type", "application/octet-stream");

      const remainingBytes = this.file.size - this.currentChunkFinalByte;
      debugger;

      if (remainingBytes - this.chunkSize < this.minimumChunk) {
        this.currentChunkFinalByte += remainingBytes;
      }

      const chunk = this.file.slice(
        this.currentChunkStartByte,
        this.currentChunkFinalByte
      );

      this.$axios.setHeader(
        "Content-Range",
        `bytes ${this.currentChunkStartByte}-${this.currentChunkFinalByte}/${this.file.size}`
      );

      console.log(this.uploadChunkSize);

      bodyFormData.append("chunk", chunk, this.file.name);
      bodyFormData.append("partNumber", this.partNumber);
      await this.$axios.$post(
        "https://upload.feblic.com/upload/part",
        bodyFormData
      );

      if (this.currentChunkFinalByte === this.file.size) {
        const res = await this.$axios.$post(
          "https://upload.feblic.com/upload/finish"
        );
        this.message = `Yay, upload completed! key is  ${res.key} and url is ${res.location}`;
        this.fileUploadStart = false;
        return;
      } else if (remainingBytes < this.chunkSize) {
        this.currentChunkStartByte = this.currentChunkFinalByte;
        this.currentChunkFinalByte =
          this.currentChunkStartByte + remainingBytes;
      } else {
        this.currentChunkStartByte = this.currentChunkFinalByte;
        this.currentChunkFinalByte =
          this.currentChunkStartByte + this.chunkSize;
      }

      this.partNumber += 1;

      this.uploadChunk();
    },
  },
};
</script>
