/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Contract } = require('fabric-contract-api');

class OTP extends Contract {

    async initLedger(ctx) {
        console.info('============= START : Initialize Ledger ===========');
        const devices = [];

        for (let i = 0; i < devices.length; i++) {
            devices[i].docType = 'device';
            await ctx.stub.putState('DEVICE' + i, Buffer.from(JSON.stringify(devices[i])));
            console.info('Added <--> ', devices[i]);
        }
        console.info('============= END : Initialize Ledger ===========');
    }

    async queryDevice(ctx, deviceNumber) {
        const deviceAsBytes = await ctx.stub.getState(deviceNumber); // get the device from chaincode state
        if (!deviceAsBytes || deviceAsBytes.length === 0) {
            throw new Error(`${deviceNumber} does not exist`);
        }
        console.log(deviceAsBytes.toString());
        return deviceAsBytes.toString();
    }

    async createDevice(ctx, deviceNumber, make, model, color, owner) {
        console.info('============= START : Create Device ===========');

        const device = {
            id,
            timestamp,
            owner,
        };

        await ctx.stub.putState(deviceNumber, Buffer.from(JSON.stringify(device)));
        console.info('============= END : Create Device ===========');
    }

    async queryAllDevices(ctx) {
        const startKey = '';
        const endKey = '';
        const allResults = [];
        for await (const {key, value} of ctx.stub.getStateByRange(startKey, endKey)) {
            const strValue = Buffer.from(value).toString('utf8');
            let record;
            try {
                record = JSON.parse(strValue);
            } catch (err) {
                console.log(err);
                record = strValue;
            }
            allResults.push({ Key: key, Record: record });
        }
        console.info(allResults);
        return JSON.stringify(allResults);
    }

    async changeDeviceOwner(ctx, deviceNumber, newOwner) {
        console.info('============= START : changeDeviceOwner ===========');

        const deviceAsBytes = await ctx.stub.getState(deviceNumber); // get the device from chaincode state
        if (!deviceAsBytes || deviceAsBytes.length === 0) {
            throw new Error(`${deviceNumber} does not exist`);
        }
        const device = JSON.parse(deviceAsBytes.toString());
        device.owner = newOwner;

        await ctx.stub.putState(deviceNumber, Buffer.from(JSON.stringify(device)));
        console.info('============= END : changeDeviceOwner ===========');
    }

}

module.exports = OTP;
