package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type PubSubAwsStackProps struct {
	awscdk.StackProps
}

func NewPubSubAwsStack(scope constructs.Construct, id string, props *PubSubAwsStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Create an sqs queue
	queue := awssqs.NewQueue(stack, jsii.String("Queue"), &awssqs.QueueProps{
		QueueName: jsii.String("test-queue"),
	})

	// Create the event bridge bus
	bus := awsevents.NewEventBus(stack, jsii.String("EventBus"), &awsevents.EventBusProps{
		EventBusName: jsii.String("test-event-bus"),
	})

	// Create the event bus rule
	rule := awsevents.NewRule(stack, jsii.String("EventBusRule"), &awsevents.RuleProps{
		EventBus: bus,
	})

	rule.AddEventPattern(&awsevents.EventPattern{
		Source:     jsii.Strings("MyCdkApplication"),
		DetailType: jsii.Strings("MessageForQueue"),
	})

	rule.AddTarget(awseventstargets.NewSqsQueue(queue, &awseventstargets.SqsQueueProps{}))

	// Emit the queue url as output
	awscdk.NewCfnOutput(stack, jsii.String("QueueUrl"), &awscdk.CfnOutputProps{
		Description: jsii.String("Url of Sqs queue"),
		Value:       queue.QueueUrl(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewPubSubAwsStack(app, "PubSubAwsStack", &PubSubAwsStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {

	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
