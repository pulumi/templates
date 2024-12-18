const { Output } = require("@pulumi/pulumi")

expect.extend({
    async toEqualOutputOf(actual, expected) {
        if (!actual instanceof Output) {
            throw new Error(`Actual value must be an Output, got ${typeof actual}`)
        }
        return new Promise(resolve =>
            actual.apply(unwrapped => {
                if (this.equals(unwrapped, expected)) {
                    resolve({
                        pass: true,
                    })
                } else {
                    resolve({
                        message: () => `expected ${this.utils.printExpected(expected)} to equal ${this.utils.printReceived(unwrapped)}`,
                        pass: false,
                    })
                }
            })
        )
    },

    async apply(actual, applyFn) {
        if (!actual instanceof Output) {
            throw new Error(`Actual value must be an Output, got ${typeof actual}`)
        }
        return new Promise(resolve =>
            actual.apply(async (...args) => {
                try {
                    await applyFn(...args)
                    resolve({
                        pass: true,
                    })
                } catch (e) {
                    resolve({
                        message: () => e.message,
                        pass: false,
                    })
                }
            })
        )
    },

})
