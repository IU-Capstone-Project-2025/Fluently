//
//  MicrophoneUsage.swift
//  Fluently
//
//  Created by Савва Пономарев on 10.07.2025.
//

import Foundation
import AVFoundation


final class MicrophoneUsage {
    var audioRecorder: AVAudioRecorder!
    var audioPlayer : AVAudioPlayer!
    var meterTimer: Timer!
    var isAudioRecordingGranted: Bool!
    var isRecording = false
    var isPlaying = false

    func startRecoding() {
        if(isRecording) {
            finishAudioRecording(success: true)
            isRecording = false
        } else {
            setupRecording()

            audioRecorder.record()
            meterTimer = Timer.scheduledTimer(timeInterval: 0.1, target:self, selector:#selector(self.updateAudioMeter(timer:)), userInfo:nil, repeats:true)
            isRecording = true
        }
    }

    @objc func updateAudioMeter(timer: Timer) {
        if audioRecorder.isRecording {
            let hr = Int((audioRecorder.currentTime / 60) / 60)
            let min = Int(audioRecorder.currentTime / 60)
            let sec = Int(audioRecorder.currentTime.truncatingRemainder(dividingBy: 60))
            let totalTimeString = String(format: "%02d:%02d:%02d", hr, min, sec)
            audioRecorder.updateMeters()
        }
    }

    func finishAudioRecording(success: Bool) {
        if success {
            audioRecorder.stop()
            audioRecorder = nil
            meterTimer.invalidate()
            print("recorded successfully.")
        } else {
            print("error")
        }
    }

    func preparePlay() {
        do {
            audioPlayer = try AVAudioPlayer(contentsOf: getFileUrl())
            audioPlayer.prepareToPlay()
        } catch {
            print(error)
        }
    }

    private func checkRecordPermission() {
        switch AVAudioSession.sharedInstance().recordPermission {
            case .undetermined:
                AVAudioSession.sharedInstance().requestRecordPermission({ (allowed) in
                    if allowed {
                        self.isAudioRecordingGranted = true
                    } else {
                        self.isAudioRecordingGranted = false
                    }
                })
                break
            case .denied:
                isAudioRecordingGranted = false
                break
            case .granted:
                isAudioRecordingGranted = true
                break
            @unknown default:
                break
        }
    }

    private func getDocumentsDirectory() -> URL {
        let paths = FileManager.default.urls(for: .documentDirectory, in: .userDomainMask)
        let documentsDirectory = paths[0]
        return documentsDirectory
    }

    private func getFileUrl() -> URL {
        let filename = "myRecording.m4a"
        let filePath = getDocumentsDirectory().appendingPathComponent(filename)
        return filePath
    }

    private func setupRecording() {
        if isAudioRecordingGranted {
            let session = AVAudioSession.sharedInstance()
            do {
                try session.setCategory(AVAudioSession.Category.playAndRecord, options: .defaultToSpeaker)
                try session.setActive(true)
                let settings = [
                    AVFormatIDKey: Int(kAudioFormatMPEG4AAC),
                    AVSampleRateKey: 44100,
                    AVNumberOfChannelsKey: 2,
                    AVEncoderAudioQualityKey:AVAudioQuality.high.rawValue
                ]
                audioRecorder = try AVAudioRecorder(url: getFileUrl(), settings: settings)
                audioRecorder.isMeteringEnabled = true
                audioRecorder.prepareToRecord()
            } catch {
                print(error)
            }
        }
    }
}

